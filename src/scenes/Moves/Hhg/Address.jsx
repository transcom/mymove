import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const formName = 'hhg_address';
const AddressWizardForm = reduxifyWizardForm(formName);

export class HHGAddress extends Component {
  handleSubmit = () => {
    const newAddress = { residential_address: this.props.values };
    // Create new HHG pickup address
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentServiceMember,
    } = this.props;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = get(currentServiceMember, 'residential_address');
    const serviceMemberId = this.props.match.params.serviceMemberId;
    return (
      <AddressWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
        initialValues={initialValues}
        additionalParams={{ serviceMemberId }}
      >
        <div className="usa-grid">
          <h3>Shipment 1 (HHG)</h3>
          <h3 className="instruction-heading">
            Now let's review your pickup and delivery locations
          </h3>

          <h4>Pickup Location</h4>

          <div className="usa-width-one-whole">
            <SwaggerField
              fieldName="street_address_1"
              swagger={this.props.schema}
              required
            />
            <SwaggerField
              fieldName="street_address_2"
              swagger={this.props.schema}
            />
            <SwaggerField
              className="usa-width-one-fourth"
              fieldName="city"
              swagger={this.props.schema}
              required
            />
            <SwaggerField
              className="usa-width-one-sixth"
              fieldName="state"
              swagger={this.props.schema}
              required
            />
            <SwaggerField
              className="usa-width-one-fourth"
              fieldName="postal_code"
              swagger={this.props.schema}
              required
            />
            <SwaggerField
              fieldName="has_secondary_pickup_address"
              swagger={this.props.schema}
              component={YesNoBoolean}
            />
            {get(
              this.props,
              'formValues.has_secondary_pickup_address',
              false,
            ) && (
              <Fragment>
                <SwaggerField
                  fieldName="street_address_1"
                  swagger={this.props.schema}
                  required
                />
                <SwaggerField
                  fieldName="street_address_2"
                  swagger={this.props.schema}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="city"
                  swagger={this.props.schema}
                  required
                />
                <SwaggerField
                  className="usa-width-one-sixth"
                  fieldName="state"
                  swagger={this.props.schema}
                  required
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="postal_code"
                  swagger={this.props.schema}
                  required
                />
              </Fragment>
            )}
            <SwaggerField
              fieldName="has_delivery_address"
              swagger={this.props.schema}
              component={YesNoBoolean}
            />
            {get(this.props, 'formValues.has_delivery_address', false) && (
              <Fragment>
                <SwaggerField
                  fieldName="street_address_1"
                  swagger={this.props.schema}
                  required
                />
                <SwaggerField
                  fieldName="street_address_2"
                  swagger={this.props.schema}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="city"
                  swagger={this.props.schema}
                  required
                />
                <SwaggerField
                  className="usa-width-one-sixth"
                  fieldName="state"
                  swagger={this.props.schema}
                  required
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="postal_code"
                  swagger={this.props.schema}
                  required
                />
              </Fragment>
            )}
            <SwaggerField
              fieldName="has_partial_sit_delivery_address"
              swagger={this.props.schema}
              component={YesNoBoolean}
            />
            {get(
              this.props,
              'formValues.has_partial_sit_delivery_address',
              false,
            ) && (
              <Fragment>
                <SwaggerField
                  fieldName="street_address_1"
                  swagger={this.props.schema}
                  required
                />
                <SwaggerField
                  fieldName="street_address_2"
                  swagger={this.props.schema}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="city"
                  swagger={this.props.schema}
                  required
                />
                <SwaggerField
                  className="usa-width-one-sixth"
                  fieldName="state"
                  swagger={this.props.schema}
                  required
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="postal_code"
                  swagger={this.props.schema}
                  required
                />
              </Fragment>
            )}
          </div>
        </div>
      </AddressWizardForm>
    );
  }
}
HHGAddress.propTypes = {
  schema: PropTypes.object.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
function mapStateToProps(state) {
  return {
    schema: get(state, 'swagger.spec.definitions.Address', {}),
    values: getFormValues(formName)(state),
    ...state.serviceMember,
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(HHGAddress);
