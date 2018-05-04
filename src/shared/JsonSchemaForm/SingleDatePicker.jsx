import React, { PureComponent, Fragment } from 'react';
import DatePicker from 'react-datepicker';

import 'react-datepicker/dist/react-datepicker.css';

export default class Example extends PureComponent {
  render() {
    const {
      meta: { error, touched },
      input: { value = null, onChange },
    } = this.props;
    return (
      <Fragment>
        <DatePicker selected={value} onChange={onChange} />
        {error && touched && <span>{error}</span>}
      </Fragment>
    );
  }
}

// I'd rather use this, but having styling trouble
// import 'react-dates/initialize';
// import { SingleDatePicker } from 'react-dates';
// import 'react-dates/lib/css/_datepicker.css';

// class SingleDatePickerField extends PureComponent {
//   state = { focused: null };
//   handleFocusChange = ({ focused }) => this.setState({ focused });

//   render() {
//     const {
//       meta: { error, touched },
//       input: { value = null, onChange },
//     } = this.props;
//     const { focused = null } = this.state;

//     return (
//       <div>
//         <SingleDatePicker
//           date={value}
//           onDateChange={onChange}
//           focused={focused}
//           onFocusChange={this.handleFocusChange}
//           id="date"
//           placeholder="date"
//           numberOfMonths={1}
//           showDefaultInputIcon
//           inputIconPosition="after"
//         />
//         {error && touched && <span>{error}</span>}
//       </div>
//     );
//   }
// }
