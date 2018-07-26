import PropTypes from 'prop-types';
import React, { Component } from 'react';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

class WeightEstimates extends Component {
  render() {
    return (
      <div className="form-section">
        <h3 className="instruction-heading">
          Enter the weight of your stuff here if you already know it
        </h3>
        <div className="usa-grid">
          <div className="usa-width-one-whole">
            <SwaggerField
              fieldName="weight_estimate"
              swagger={this.props.schema}
            />
          </div>
        </div>
        <h3 className="instruction-heading">
          Now enter the weight of your Pro-Gear
        </h3>
        <p>
          You are entitled to move up to 2000 lbs. of pro-gear and 500 lbs. of
          spouse pro-gear. Pro-Gear includes uniforms, deployment gear, and any
          other gear you or your spouse need to perform your jobs.
        </p>
        <div className="usa-grid">
          <div className="usa-width-one-whole">
            <SwaggerField
              fieldName="progear_weight_estimate"
              swagger={this.props.schema}
            />
            <SwaggerField
              fieldName="spouse_progear_weight_estimate"
              swagger={this.props.schema}
            />
          </div>
        </div>
      </div>
    );
  }
}

WeightEstimates.propTypes = {
  schema: PropTypes.object.isRequired,
  error: PropTypes.object,
  formValues: PropTypes.object,
};

export default WeightEstimates;
