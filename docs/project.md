# Food Delivery

> This is the authoritative student-facing project specification and API contract.
> Changes are recorded in the [specification changelog](changelog.md).

Build two services for a small food-delivery system:

- **Order Service** manages carts, orders, restaurant decisions, and delivery requests.
- **Restaurant Service** manages restaurants, menus, prices, and availability.

The project excludes geographic search, couriers, live location, payment, promotions, and delivery-time estimates. Customer and restaurant identifiers are passed in headers instead of using real authentication.

## Business rules

### Restaurant and Menu

- A Restaurant may accept or pause new Orders.
- A Restaurant has one currency and a minimum Order amount.
- A Menu Item belongs to one Restaurant, uses its currency, and may be available or unavailable.
- Restaurant and Menu Item identifiers are globally unique strings.
- Restaurant reads are public. Only the matching Restaurant operator may create or replace its Restaurant and Menu Items.

### Cart

- A Customer has at most one active Cart.
- A Cart may contain Items from only one Restaurant. An empty Cart has no Restaurant.
- An Item must exist and be available when it is added.
- A Cart has at most one line per Item, 20 lines total, and 50 total Item units.
- Quantity is an absolute value from 1 through 10. Repeating the same request has the same result.
- Preparation instructions are optional and at most 250 characters.
- Adding an Item from another Restaurant fails without changing the Cart.
- Removing an absent line succeeds. Removing the final line clears the Restaurant association.
- Cart prices are estimates; checkout always uses current Restaurant Service prices.

### Order

- An Order is created from the complete active Cart.
- The client supplies only a delivery address: non-empty after trimming and at most 500 characters.
- Before creation, Order Service revalidates the Restaurant, Menu Items, availability, prices, currency, and minimum Order amount.
- Order Lines snapshot Item names, prices, quantities, and preparation instructions.
- A new Order is `PENDING`. Success clears the Cart; failure changes nothing.
- Only the owning Customer may read or cancel an Order.
- Only the owning Restaurant may read or advance an Order.
- A Customer may cancel only a `PENDING` Order.
- A Restaurant may accept or reject only a `PENDING` Order. Rejection requires a reason.
- Accepting creates exactly one Delivery in `WAITING_FOR_PREPARATION`, idempotently by `order_id`.
- If Delivery creation fails, the Order remains `PENDING`.
- Marking an Order ready changes its Delivery to `IN_DELIVERY`.
- If that Delivery transition fails, the Order remains `PREPARING`.
- Repeating an already completed command succeeds without changing timestamps. Other skipped, backward, or conflicting transitions fail.

## Status models

```text
Order:
                    ┌──────────→ REJECTED
                    │
PENDING ────────────┼──────────→ CANCELLED
                    │
                    └→ ACCEPTED → PREPARING → READY_FOR_PICKUP

Delivery:
WAITING_FOR_PREPARATION → IN_DELIVERY
```

`IN_DELIVERY` is terminal in this project. Courier assignment and delivery confirmation are not modeled.

> **Review note for Lesson 3:** These are two distinct moments. When a Restaurant accepts a `PENDING` Order, the Order becomes `ACCEPTED` and a Delivery is created in `WAITING_FOR_PREPARATION`. When the Restaurant later marks a `PREPARING` Order ready, the Order becomes `READY_FOR_PICKUP` and that existing Delivery moves to `IN_DELIVERY`. Please review these transitions before publishing Lesson 3; this note does not add a new status or change the rules above.

## API conventions

- JSON request and response bodies.
- Money is an integer number of minor currency units.
- Times are RFC 3339 in UTC.
- Order Service Customer endpoints require `X-Customer-ID`.
- Order Service Restaurant endpoints and Restaurant Service `PUT` endpoints require `X-Restaurant-ID`.
- Restaurant Service `GET` endpoints do not require an identity header.
- `X-Restaurant-ID` must match the Restaurant identified by the path or Order.
- Unknown JSON fields, malformed bodies, and missing identity headers are invalid requests.

Error envelope:

```json
{
  "error": {
    "code": "CART_RESTAURANT_CONFLICT",
    "message": "a cart can contain items from only one restaurant",
    "details": {}
  }
}
```

HTTP status groups: `400` invalid request, `403` forbidden, `404` not found, `409` conflict or invalid transition, `422` business-rule violation, `503` dependency unavailable, and `500` unexpected failure.

## Restaurant Service contract

### Representations

```text
Restaurant {
  id string
  name string
  accepting_orders bool
  minimum_order_minor int64
  currency string
}

MenuItem {
  id string
  name string
  description string
  price_minor int64
  currency string
  available bool
}
```

### Operations

| Method and path                                        | Request                   | Success response                           |
| ------------------------------------------------------ | ------------------------- | ------------------------------------------ |
| `PUT /restaurants/{restaurant_id}`                     | `Restaurant` without `id` | `200` with `Restaurant`                    |
| `GET /restaurants`                                     | —                         | `200 { restaurants: Restaurant[] }`        |
| `GET /restaurants/{restaurant_id}`                     | —                         | `200` with `Restaurant`                    |
| `PUT /restaurants/{restaurant_id}/menu/{menu_item_id}` | `MenuItem` without `id`   | `200` with `MenuItem`                      |
| `GET /restaurants/{restaurant_id}/menu`                | —                         | `200 { restaurant_id, items: MenuItem[] }` |
| `GET /restaurants/{restaurant_id}/menu/{menu_item_id}` | —                         | `200` with `MenuItem`                      |

Both `PUT` operations create or replace the resource and are idempotent.

Errors: `RESTAURANT_NOT_FOUND`, `RESTAURANT_ACCESS_DENIED`, `MENU_ITEM_NOT_FOUND`, `INVALID_RESTAURANT`, `INVALID_MENU_ITEM`, `CURRENCY_MISMATCH`.

## Order Service contract

### Representations

```text
CartItem {
  menu_item_id string
  name string
  unit_price_minor int64
  currency string
  quantity int32
  instructions string
}

Cart {
  restaurant_id string|null
  items CartItem[]
  subtotal_minor int64
  currency string|null
}

OrderItem {
  menu_item_id string
  name string
  unit_price_minor int64
  quantity int32
  instructions string
}

Order {
  id string
  customer_id string
  restaurant_id string
  status PENDING|ACCEPTED|PREPARING|READY_FOR_PICKUP|REJECTED|CANCELLED
  items OrderItem[]
  subtotal_minor int64
  currency string
  delivery_address string
  rejection_reason string|null
  delivery_status WAITING_FOR_PREPARATION|IN_DELIVERY|null
  created_at timestamp
  updated_at timestamp
}
```

An absent Cart is returned as `200` with a null Restaurant and currency, an empty Item list, and zero subtotal. Delivery status is null until acceptance and stays null for rejected or cancelled Orders.

### Customer operations

| Method and path                     | Request                                     | Success response           |
| ----------------------------------- | ------------------------------------------- | -------------------------- |
| `GET /cart`                         | —                                           | `200` with `Cart`          |
| `PUT /cart/items/{menu_item_id}`    | `{ restaurant_id, quantity, instructions }` | `200` with updated `Cart`  |
| `DELETE /cart/items/{menu_item_id}` | —                                           | `204`                      |
| `POST /orders`                      | `{ delivery_address }`                      | `201` with `Order`         |
| `GET /orders/{order_id}`            | —                                           | `200` with `Order`         |
| `GET /orders`                       | —                                           | `200 { orders: Order[] }`  |
| `POST /orders/{order_id}/cancel`    | —                                           | `200` with updated `Order` |

### Restaurant operations

| Method and path                                        | Request      | Success response           |
| ------------------------------------------------------ | ------------ | -------------------------- |
| `GET /restaurant/orders`                               | —            | `200 { orders: Order[] }`  |
| `POST /restaurant/orders/{order_id}/accept`            | —            | `200` with updated `Order` |
| `POST /restaurant/orders/{order_id}/reject`            | `{ reason }` | `200` with updated `Order` |
| `POST /restaurant/orders/{order_id}/start-preparation` | —            | `200` with updated `Order` |
| `POST /restaurant/orders/{order_id}/ready`             | —            | `200` with updated `Order` |

Errors: `CUSTOMER_NOT_FOUND`, `RESTAURANT_NOT_FOUND`, `RESTAURANT_NOT_ACCEPTING_ORDERS`, `MENU_ITEM_NOT_FOUND`, `MENU_ITEM_UNAVAILABLE`, `CART_RESTAURANT_CONFLICT`, `CART_LIMIT_EXCEEDED`, `INVALID_QUANTITY`, `INVALID_INSTRUCTIONS`, `EMPTY_CART`, `MINIMUM_ORDER_NOT_REACHED`, `INVALID_DELIVERY_ADDRESS`, `ORDER_NOT_FOUND`, `ORDER_ACCESS_DENIED`, `REJECTION_REASON_REQUIRED`, `INVALID_ORDER_TRANSITION`, `DELIVERY_CREATION_FAILED`, `DELIVERY_TRANSITION_FAILED`, `DEPENDENCY_UNAVAILABLE`.
