import PropTypes from 'prop-types';
import React, { Component } from 'react';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

class ProgearEstimates extends Component {
  render() {
    return (
      <div className="form-section">
        <h3 className="instruction-heading">Now enter the weight of your Pro-Gear</h3>
        <p>
          Pro-Gear includes uniforms, deployment gear, and any other gear you or your spouse need to perform your jobs.
        </p>
        <div className="usa-grid">
          <div className="usa-width-one-whole">
            <SwaggerField fieldName="progear_weight_estimate" swagger={this.props.schema} />
            <SwaggerField fieldName="spouse_progear_weight_estimate" swagger={this.props.schema} />
          </div>
        </div>
      </div>
    );
  }
}

ProgearEstimates.propTypes = {
  schema: PropTypes.object.isRequired,
  formValues: PropTypes.object,
};

export default ProgearEstimates;
