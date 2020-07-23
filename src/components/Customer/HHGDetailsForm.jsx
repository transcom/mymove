import React, { Component } from 'react';
import { arrayOf, string, bool, shape, func } from 'prop-types';
import { connect } from 'react-redux';
import { Formik } from 'formik';
import { Fieldset, Radio, Label } from '@trussworks/react-uswds';
import { push } from 'connected-react-router';

import { Form } from '../form/Form';
import { DatePickerInput, TextInput } from '../form/fields';
import { AddressFields } from '../form/AddressFields/AddressFields';
import { ContactInfoFields } from '../form/ContactInfoFields/ContactInfoFields';

import { createMTOShipment as createMTOShipmentAction } from 'shared/Entities/modules/mtoShipments';
import { WizardPage } from 'shared/WizardPage';
import { MTOAgentType } from 'shared/constants';
import { formatSwaggerDate } from 'shared/formatters';

class HHGDetailsForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hasDeliveryAddress: false,
    };
  }

  // TODO: when we can pull in initialValues from redux, set state.hasDeliveryAddress to true if a delivery address exists

  handleChangeHasDeliveryAddress = () => {
    this.setState((prevState) => {
      return { hasDeliveryAddress: !prevState.hasDeliveryAddress };
    });
  };

  submitMTOShipment = ({
    requestedPickupDate,
    requestedDeliveryDate,
    pickupLocation,
    deliveryLocation,
    receivingAgent,
    releasingAgent,
    remarks,
  }) => {
    const { createMTOShipment } = this.props;
    const { hasDeliveryAddress } = this.state;
    const mtoShipment = {
      // TODO: Use moveTaskOrderID when it is available
      moveTaskOrderID: '5d4b25bb-eb04-4c03-9a81-ee0398cb779e',
      shipmentType: 'HHG',
      requestedPickupDate: formatSwaggerDate(requestedPickupDate),
      requestedDeliveryDate: formatSwaggerDate(requestedDeliveryDate),
      customerRemarks: remarks,
      pickupAddress: {
        street_address_1: pickupLocation.mailingAddress1,
        street_address_2: pickupLocation.mailingAddress2,
        city: pickupLocation.city,
        state: pickupLocation.state,
        postal_code: pickupLocation.zip,
        country: pickupLocation.country,
      },
      agents: [
        { ...releasingAgent, agentType: MTOAgentType.RELEASING },
        { ...receivingAgent, agentType: MTOAgentType.RECEIVING },
      ],
    };
    if (hasDeliveryAddress) {
      mtoShipment.destinationAddress = {
        street_address_1: deliveryLocation.mailingAddress1,
        street_address_2: deliveryLocation.mailingAddress2,
        city: deliveryLocation.city,
        state: deliveryLocation.state,
        postal_code: deliveryLocation.zip,
        country: deliveryLocation.country,
      };
    }
    createMTOShipment(mtoShipment);
  };

  render() {
    // TODO: replace minimal styling with actual styling during UI phase
    const { initialValues, pageKey, pageList, match } = this.props;
    const { hasDeliveryAddress } = this.state;
    const fieldsetClasses = 'margin-top-2';
    return (
      <Formik initialValues={initialValues}>
        {({ handleChange, values }) => (
          <WizardPage
            match={match}
            push={push}
            pageKey={pageKey}
            pageList={pageList}
            handleSubmit={() => this.submitMTOShipment(values)}
          >
            <Form>
              <Fieldset legend="Pickup date" className={fieldsetClasses}>
                <DatePickerInput
                  name="requestedPickupDate"
                  label="Requested pickup date"
                  id="requestedPickupDate"
                  value={values.requestedPickupDate}
                />
              </Fieldset>
              <span className="usa-hint" id="pickupDateHint">
                Your movers will confirm this date or one shortly before or after.
              </span>
              <AddressFields
                name="pickupLocation"
                legend="Pickup location"
                className={fieldsetClasses}
                handleChange={handleChange}
                values={values.pickupLocation}
              />
              <ContactInfoFields
                name="releasingAgent"
                legend="Releasing agent"
                className={fieldsetClasses}
                subtitle="Who can allow the movers to take your stuff if you're not there?"
                handleChange={handleChange}
                values={values.releasingAgent}
              />
              <Fieldset legend="Delivery date" className={fieldsetClasses}>
                <DatePickerInput
                  name="requestedDeliveryDate"
                  label="Requested delivery date"
                  id="requestedDeliveryDate"
                  value={values.requestedDeliveryDate}
                />
                <span className="usa-hint" id="deliveryDateHint">
                  Your movers will confirm this date or one shortly before or after.
                </span>
              </Fieldset>
              <Fieldset legend="Delivery location" className={fieldsetClasses}>
                <Label>Do you know your delivery address?</Label>
                <div className="display-flex margin-top-1">
                  <Radio
                    id="has-delivery-address"
                    label="Yes"
                    name="hasDeliveryAddress"
                    onChange={this.handleChangeHasDeliveryAddress}
                    checked={hasDeliveryAddress}
                  />
                  <Radio
                    id="no-delivery-address"
                    label="No"
                    name="hasDeliveryAddress"
                    checked={!hasDeliveryAddress}
                    onChange={this.handleChangeHasDeliveryAddress}
                  />
                </div>
                {hasDeliveryAddress ? (
                  <AddressFields
                    name="deliveryLocation"
                    className={fieldsetClasses}
                    handleChange={handleChange}
                    values={values.deliveryLocation}
                  />
                ) : (
                  <>
                    <div className={fieldsetClasses}>We can use the zip of your new duty station.</div>
                    <div>[City], [State] [New duty station zip]</div>
                  </>
                )}
              </Fieldset>
              <ContactInfoFields
                name="receivingAgent"
                legend="Receiving agent"
                className={fieldsetClasses}
                subtitle="Who can take delivery for you if the movers arrive and you're not there?"
                handleChange={handleChange}
                values={values.receivingAgent}
              />
              <Fieldset legend="Remarks" className={fieldsetClasses}>
                <Label hint="(optional)">Anything else you would like us to know?</Label>
                <TextInput
                  name="remarks"
                  id="remarks"
                  maxLength={1500}
                  onChange={handleChange}
                  value={values.remarks}
                />
              </Fieldset>
            </Form>
          </WizardPage>
        )}
      </Formik>
    );
  }
}

HHGDetailsForm.propTypes = {
  pageKey: string.isRequired,
  pageList: arrayOf(string).isRequired,
  initialValues: shape({
    requestedPickupDate: string,
    pickupLocation: shape({
      mailingAddress1: string,
      mailingAddress2: string,
      city: string,
      state: string,
      zip: string,
    }),
    requestedDeliveryDate: string,
    deliveryLocation: shape({
      mailingAddress1: string,
      mailingAddress2: string,
      city: string,
      state: string,
      zip: string,
    }),
    releasingAgent: shape({
      firstName: string,
      lastName: string,
      phone: string,
      email: string,
    }),
    receivingAgent: shape({
      firstName: string,
      lastName: string,
      phone: string,
      email: string,
    }),
    remarks: string,
  }),
  match: shape({
    isExact: bool.isRequired,
    params: shape({
      moveId: string.isRequired,
    }),
    path: string.isRequired,
    url: string.isRequired,
  }).isRequired,
  createMTOShipment: func.isRequired,
};

HHGDetailsForm.defaultProps = {
  initialValues: {},
};

const mapDispatchToProps = {
  createMTOShipment: createMTOShipmentAction,
};

export default connect(null, mapDispatchToProps)(HHGDetailsForm);
