# Deprecate use of lodash over time

We are currently using the utility library [lodash](https://lodash.com/) throughout the frontend codebase. However, many of the methods provided by lodash have been superseded by native ES6+ methods. Even the latest ECMAScript functionality that is not natively supported by our target browsers (such as IE11) will be transpiled and/or polyfilled with our existing Webpack & Babel configuration. Continuing to use lodash methods as-is has two negative effects:

- Inconsistent usage as some contributors will opt to use native JS functions instead of the lodash alternative
- An additional 3rd party dependency that we need to keep updated, and adds to our bundle size

I am proposing that we fully deprecate usage of lodash methods that have native JavaScript equivalents. We can use the [You Don’t Need Lodash/Underscore](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore) ESLint plugin to begin enforcing consistent usage on a per-method basis. Additionally, we can use the [lodash-webpack-plugin](https://www.npmjs.com/package/lodash-webpack-plugin) and [babel-plugin-lodash](https://www.npmjs.com/package/babel-plugin-lodash) to cherry-pick lodash methods with no native equivalent that we may want to continue using (such as `uniqueId`), and result in a smaller bundle.

Additionally, this strategy leaves us open to completely removing lodash as a dependency in the future if we so choose. I'm focusing on bundle size over individual method benchmarks, because it is [well-documented](https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e) that reducing parse & compile time is the most effective way to improve performance on mobile devices.

_Note that some of the native equivalents are not exactly the same as the lodash versions. The links below point to documentation for any differences._

**Methods to deprecate (links to eslint plugin docs):**

- [bind](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_bind)
- [concat](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_concat)
- [debounce](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_debounce)
- [every](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_every)
- [filter](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_filter)
- [forEach](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_each)
- [get](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_get)
- [head](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_head-and-_tail)
- [includes](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_includes)
- [isEmpty](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_isempty)
- [isFinite](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_isfinite)
- [isInteger](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_isInteger)
- [isNil](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_isnil)
- [isNull](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_isnull)
- [last](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_last)
- [map](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_map)
- [omit](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_omit)
- [pick](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_pick)
- [reject](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_reject)
- [some](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_some)
- [sortBy](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_sortby-and-_orderby)
- [startsWith](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_startsWith)
- [toUpper](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_toupper)
- [uniq](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_uniq)
- [without](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore#_without)

**Methods to retain (links to lodash docs):**

- [capitalize](https://lodash.com/docs/4.17.15#capitalize)
- [clone](https://lodash.com/docs/4.17.15#clone)
- [cloneDeep](https://lodash.com/docs/4.17.15#cloneDeep)
- [findKey](https://lodash.com/docs/4.17.15#findKey)
- [findLast](https://lodash.com/docs/4.17.15#findLast)
- [mapValues](https://lodash.com/docs/4.17.15#mapValues)
- [memoize](https://lodash.com/docs/4.17.15#memoize)
- [snakeCase](https://lodash.com/docs/4.17.15#snakeCase)
- [startCase](https://lodash.com/docs/4.17.15#startCase)
- [sum](https://lodash.com/docs/4.17.15#sum)
- [union](https://lodash.com/docs/4.17.15#union)
- [uniqueId](https://lodash.com/docs/4.17.15#uniqueId)

## Considered Alternatives

- Deprecate lodash methods that have close native equivalents
- Deprecate lodash entirely
- Switch from native methods to only lodash methods
- Don't change lodash usage, but add Babel & Webpack plugins
- Do nothing

## Decision Outcome

- Chosen Alternative: Deprecate lodash methods that have close native equivalents.
- Additionally, set up the [lodash-webpack-plugin](https://www.npmjs.com/package/lodash-webpack-plugin) and [babel-plugin-lodash](https://www.npmjs.com/package/babel-plugin-lodash) to optimize how lodash adds to our bundle size.
- `+` Enforce one consistent style throughout the codebase
- `+` Can replace existing lodash methods on a per-method basis using the [ESLint plugin](https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore) (and prevent more instances from being added)
- `+` Lodash imports will be optimized and should reduce bundle size
- `-` Need to socialize and educate engineers on the native methods to use instead

## Pros and Cons of the Alternatives

### Deprecate lodash entirely

- `+` Removes a significant library dependency
- `+` Enforce one consistent style throughout the codebase
- `-` Would have to immediately replace existing lodash methods throughout the code
- `-` Need to find alternate solutions for lodash methods that have no native equivalent (such as `uniqueId`)

### Switch from native methods to only lodash methods

- `+` Enforce one consistent style throughout the codebase
- `-` Increased reliance on a third-party library, and learning curve for engineers who are more familiar with native JS than with lodash
- `-` Our code is further removed from native JS implementations and future proposals

### Don't change lodash usage, but add Babel & Webpack plugins

- `+` Requires no code or habit changes (other than Webpack/Babel config)
- `+` Lodash imports will be optimized and should reduce bundle size
- `-` Perpetuates inconsistent code styles and confusion for new engineers

### Do nothing

- `+` Requires no immediate effort
- `-` Perpetuates inconsistent code styles and confusion for new engineers
- `-` Lodash continues to account for ~98KB of our bundled frontend assets
