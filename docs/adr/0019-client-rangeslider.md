# _Range Slider React Component_

**User Story:** _[#155911084](https://www.pivotaltracker.com/story/show/155911084)_ <!-- optional -->

We need a range slider component for the PPM Incentives screen. This needs to be

* accessible
* styleable
* responsive

## Considered Alternatives

* _[react-rangeslider](https://whoisandy.github.io/react-rangeslider/)_
* _[rc-slider](http://react-component.github.io/slider/)_
* _[rheostat](https://github.com/airbnb/rheostat)_
* _[rea11y-sliders](https://github.com/felixgirault/rea11y/blob/master/docs/Slider.md)_
* _[react-html5-slider](https://github.com/brownieboy/react-html5-slider)_

## Decision Outcome

* Chosen Alternative: _react-rangeslider_
* _It was the only one that would build, was controllable by keyboard, and looked decent with USWDS CSS_

## Pros and Cons of the Alternatives <!-- optional -->

### _react-rangeslider_

* `+` _it worked_
* `+` _supported keyboard controls_
* `-` _examples require component state_

### _rc-slider_

* `+` _it worked_
* `-` _could only be controlled with mouse_

### _rheostat_

* `+` _AirBnb created_
* `+` _documentation seemed nice_
* `-` _even storybook they deliver has broken styles_
* `-` \*layout was unusable with uswds styles

### _rea11y-sliders_

* `-` _would not compile_

### _react-html5-slider_

* `-` _would not compile_
