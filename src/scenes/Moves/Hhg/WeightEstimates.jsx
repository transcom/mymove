import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
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
        <h3 className="instruction-heading usa-heading">
          Now enter the weight of your Pro-Gear
        </h3>
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

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
function mapStateToProps(state) {
  return {};
}
export default connect(mapStateToProps, mapDispatchToProps)(WeightEstimates);
