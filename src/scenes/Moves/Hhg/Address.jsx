import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, FormSection } from 'redux-form';

import {
  createOrUpdateShipment,
  selectShipment,
} from 'shared/Entities/modules/shipments';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';

import './Address.css';

const formName = 'shipment';
const AddressWizardForm = reduxifyWizardForm(formName);

export class ShipmentAddress extends Component {
  handleSubmit = () => {
    const moveId = this.props.match.params.moveId;
    const shipment = this.props.formValues;
    this.props
      .createOrUpdateShipment(moveId, shipment)
      .then(() => {
        console.log('You did it!');
      })
      .catch(err => {
        this.setState({
          shipmentError: err,
        });
      });
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      initialValues,
    } = this.props;

    const hasSecondary = get(
      this.props,
      'formValues.has_secondary_pickup_address',
      false,
    );
    const hasDelivery = get(
      this.props,
      'formValues.has_delivery_address',
      false,
    );

    return (
      <AddressWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={false}
        serverError={error}
        initialValues={initialValues}
      >
        <div className="usa-grid">
          <h3>Shipment 1 (HHG)</h3>
          <h3 className="instruction-heading">
            Now let's review your pickup and delivery locations
          </h3>

          <h4>Pickup Location</h4>

          <div className="usa-width-one-whole">
            <FormSection name="pickup_address">
              <div className="address-segment">
                <SwaggerField
                  fieldName="street_address_1"
                  swagger={this.props.schema.properties.pickup_address}
                  required
                />
                <SwaggerField
                  fieldName="street_address_2"
                  swagger={this.props.schema.properties.pickup_address}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="city"
                  swagger={this.props.schema.properties.pickup_address}
                  required
                />
                <SwaggerField
                  className="usa-width-one-sixth"
                  fieldName="state"
                  swagger={this.props.schema.properties.pickup_address}
                  required
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="postal_code"
                  swagger={this.props.schema.properties.pickup_address}
                  required
                />
              </div>
            </FormSection>
            <SwaggerField
              className="radio-title"
              fieldName="has_secondary_pickup_address"
              swagger={this.props.schema}
              component={YesNoBoolean}
            />
            {hasSecondary && (
              <Fragment>
                <FormSection name="secondaryPickupAddress">
                  <div className="address-segment">
                    <SwaggerField
                      fieldName="street_address_1"
                      swagger={
                        this.props.schema.properties.secondary_pickup_address
                      }
                      required={hasSecondary}
                    />
                    <SwaggerField
                      fieldName="street_address_2"
                      swagger={
                        this.props.schema.properties.secondary_pickup_address
                      }
                    />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="city"
                      swagger={
                        this.props.schema.properties.secondary_pickup_address
                      }
                      required={hasSecondary}
                    />
                    <SwaggerField
                      className="usa-width-one-sixth"
                      fieldName="state"
                      swagger={
                        this.props.schema.properties.secondary_pickup_address
                      }
                      required={hasSecondary}
                    />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="postal_code"
                      swagger={
                        this.props.schema.properties.secondary_pickup_address
                      }
                      required={hasSecondary}
                    />
                  </div>
                </FormSection>
              </Fragment>
            )}
            <h4>Delivery location</h4>
            <SwaggerField
              className="radio-title"
              fieldName="has_delivery_address"
              swagger={this.props.schema}
              component={YesNoBoolean}
            />
            {hasDelivery && (
              <Fragment>
                <FormSection name="deliveryAddress">
                  <div className="address-segment">
                    <SwaggerField
                      fieldName="street_address_1"
                      swagger={this.props.schema.properties.delivery_address}
                      required={hasDelivery}
                    />
                    <SwaggerField
                      fieldName="street_address_2"
                      swagger={this.props.schema.properties.delivery_address}
                    />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="city"
                      swagger={this.props.schema.properties.delivery_address}
                      required={hasDelivery}
                    />
                    <SwaggerField
                      className="usa-width-one-sixth"
                      fieldName="state"
                      swagger={this.props.schema.properties.delivery_address}
                      required={hasDelivery}
                    />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="postal_code"
                      swagger={this.props.schema.properties.delivery_address}
                      required={hasDelivery}
                    />
                  </div>
                </FormSection>
              </Fragment>
            )}
          </div>
        </div>
      </AddressWizardForm>
    );
  }
}
ShipmentAddress.propTypes = {
  schema: PropTypes.object.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { createOrUpdateShipment, selectShipment },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    schema: get(state, 'swagger.spec.definitions.Shipment', {}),
    formValues: getFormValues(formName)(state),
    move: get(state, 'moves.currentMove', {}),
    initialValues: get(state, 'moves.currentMove.shipments[0]', {}),
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(ShipmentAddress);
