# Use camelCase for API params

All API params should use the same style of casing (camel, snake, kebab, etc.)to maintain consistency across the code base.

## Considered Alternatives

* Leave everything as-is. Allow mixed cases in yaml files.
* Choose one case style and implement across the code base.

## Decision Outcome

* Chosen Alternative: **Use camelCase for all API params.**

Consistency is important and using mixed cases is confusing. We selected camelCase because it's predominant in the code base and required the least amount of effort to implement.

It's also used in `orders.yaml`, which we cannot change because of different services that download and use it. By sticking with camelCase, all of our yaml files will use the same standard and prevent further confusion.

## Pros and Cons of the Alternatives

### Leave everything as-is. Allow mixed cases in yaml files

* `+` No changes needed
* `-` Confusing for developers
* `-` No documented standard sets us up to repeat the same conversation in the future
* `-` Code base is inconsistent

### Choose one case style and implement across the code base

* `+` Clear standard for developers to use
* `+` Saves time and energy discussing preferred styles in the future
* `+` Code base is consistent
* `-` Requires manual changes in yaml files and across front-end components