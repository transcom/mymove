import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { uniqueId } from 'lodash';

const YesNoBoolean = props => {
  let value, onChange;
  if (props.input) {
    value = Boolean(props.input.value);
    onChange = props.input.onChange;
  } else {
    value = Boolean(props.value);
    onChange = props.onChange;
  }
  const yesId = uniqueId('yes_no_');
  const noId = uniqueId('yes_no_');
  const localOnChange = event => {
    onChange(event.target.value === 'yes');
  };

  return (
    <Fragment>
      <input className="inline_radio" id={yesId} type="radio" value="yes" onChange={localOnChange} checked={value} />
      <label className="inline_radio" htmlFor={yesId}>
        Yes
      </label>
      <input className="inline_radio" id={noId} value="no" type="radio" onChange={localOnChange} checked={!value} />
      <label className="inline_radio" htmlFor={noId}>
        No
      </label>
    </Fragment>
  );
};
YesNoBoolean.propTypes = {
  input: PropTypes.shape({
    value: PropTypes.oneOfType([PropTypes.string, PropTypes.bool]).isRequired,
    onChange: PropTypes.func.isRequired,
  }),
};
export default YesNoBoolean;
