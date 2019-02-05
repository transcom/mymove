import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { FormSection } from 'redux-form';

import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './Address.css';

export class ShipmentAddress extends Component {
  handleClick = () => {
    alert(
      "If you don't know your destination address before the movers arrive at your destination city, they will put your belongings into storage. Depending on the season, it can take from a few days to a month to get your stuff out of storage. If you are concerned that you might not have a new address in time, make sure to keep essential living items with you.",
    );
  };

  render() {
    const hasSecondary = get(this.props, 'formValues.has_secondary_pickup_address', false);
    const hasDelivery = get(this.props, 'formValues.has_delivery_address', false);
    const addressSchema = get(this.props, 'schema.properties.pickup_address', {});

    return (
      <div className="form-section">
        <h3 className="instruction-heading">Pickup and delivery locations</h3>

        <div className="usa-grid">
          <div className="usa-width-one-whole">
            <h4>Pickup Location</h4>
            <FormSection name="pickup_address">
              <div className="address-segment usa-grid">
                <SwaggerField fieldName="street_address_1" swagger={addressSchema} required />
                <SwaggerField fieldName="street_address_2" swagger={addressSchema} />
                <SwaggerField className="usa-width-one-fourth" fieldName="city" swagger={addressSchema} required />
                <SwaggerField className="usa-width-one-sixth" fieldName="state" swagger={addressSchema} required />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="postal_code"
                  swagger={addressSchema}
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
                <h4>Secondary Pickup Location</h4>
                <FormSection name="secondary_pickup_address">
                  <div className="address-segment usa-grid">
                    <SwaggerField fieldName="street_address_1" swagger={addressSchema} required={hasSecondary} />
                    <SwaggerField fieldName="street_address_2" swagger={addressSchema} />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="city"
                      swagger={addressSchema}
                      required={hasSecondary}
                    />
                    <SwaggerField
                      className="usa-width-one-sixth"
                      fieldName="state"
                      swagger={addressSchema}
                      required={hasSecondary}
                    />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="postal_code"
                      swagger={addressSchema}
                      required={hasSecondary}
                    />
                  </div>
                </FormSection>
              </Fragment>
            )}
            <h4>Delivery Location</h4>
            <div className="delivery-location">
              <SwaggerField
                className="radio-title"
                fieldName="has_delivery_address"
                swagger={this.props.schema}
                component={YesNoBoolean}
              />
              <div>
                <a className="tool-tip" onClick={this.handleClick}>
                  What happens if I don't know before I move?
                </a>
              </div>
            </div>
            {hasDelivery && (
              <Fragment>
                <FormSection name="delivery_address">
                  <div className="address-segment usa-grid">
                    <SwaggerField fieldName="street_address_1" swagger={addressSchema} required={hasDelivery} />
                    <SwaggerField fieldName="street_address_2" swagger={addressSchema} />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="city"
                      swagger={addressSchema}
                      required={hasDelivery}
                    />
                    <SwaggerField
                      className="usa-width-one-sixth"
                      fieldName="state"
                      swagger={addressSchema}
                      required={hasDelivery}
                    />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="postal_code"
                      swagger={addressSchema}
                      required={hasDelivery}
                    />
                  </div>
                </FormSection>
              </Fragment>
            )}
          </div>
        </div>
      </div>
    );
  }
}
ShipmentAddress.propTypes = {
  schema: PropTypes.object.isRequired,
  error: PropTypes.object,
  formValues: PropTypes.object,
};

export default ShipmentAddress;
