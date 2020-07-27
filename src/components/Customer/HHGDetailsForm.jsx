import React, { Component } from 'react';
import { arrayOf, string, bool, shape, func } from 'prop-types';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { Formik } from 'formik';
import { Fieldset, Radio, Label } from '@trussworks/react-uswds';

import { Form } from '../form/Form';
import { DatePickerInput, TextInput } from '../form/fields';
import { AddressFields } from '../form/AddressFields/AddressFields';
import { ContactInfoFields } from '../form/ContactInfoFields/ContactInfoFields';

import { createMTOShipment as createMTOShipmentAction } from 'shared/Entities/modules/mtoShipments';
import { showLoggedInUser as showLoggedInUserAction, selectLoggedInUser } from 'shared/Entities/modules/user';
import { WizardPage } from 'shared/WizardPage';
import { MTOAgentType } from 'shared/constants';
import { formatSwaggerDate } from 'shared/formatters';
import Checkbox from 'shared/Checkbox';

class HHGDetailsForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hasDeliveryAddress: false,
      useCurrentResidence: false,
      initialValues: {},
    };
  }

  componentDidMount() {
    const { showLoggedInUser } = this.props;
    showLoggedInUser();
  }

  // TODO: when we can pull in initialValues from redux, set state.hasDeliveryAddress to true if a delivery address exists
  handleChangeHasDeliveryAddress = () => {
    this.setState((prevState) => {
      return { hasDeliveryAddress: !prevState.hasDeliveryAddress };
    });
  };

  // Use current residence
  handleUseCurrentResidenceChange = (currentValues) => {
    // eslint-disable-next-line react/destructuring-assignment
    this.setState(
      (state) => ({ useCurrentResidence: !state.useCurrentResidence }),
      () => {
        const { initialValues, useCurrentResidence } = this.state;
        const { currentResidence } = this.props;
        if (useCurrentResidence) {
          this.setState({
            // eslint-disable-next-line prettier/prettier
            initialValues: {
              ...initialValues,
              ...currentValues,
              pickupLocation: {
                mailingAddress1: currentResidence.street_address_1,
                mailingAddress2: currentResidence.street_address_2,
                city: currentResidence.city,
                state: currentResidence.state,
                zip: currentResidence.postal_code,
              },
            },
          });
        } else {
          this.setState({
            initialValues: {
              ...initialValues,
              ...currentValues,
              pickupLocation: {
                mailingAddress1: '',
                mailingAddress2: '',
                city: '',
                state: '',
                zip: '',
              },
            },
          });
        }
      },
    );
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
    const { createMTOShipment, moveTaskOrderID } = this.props;
    const { hasDeliveryAddress } = this.state;
    const mtoShipment = {
      moveTaskOrderID,
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
      agents: [],
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
    if (releasingAgent) {
      mtoShipment.agents.push({ ...releasingAgent, agentType: MTOAgentType.RELEASING });
    }

    if (receivingAgent) {
      mtoShipment.agents.push({ ...receivingAgent, agentType: MTOAgentType.RECEIVING });
    }
    createMTOShipment(mtoShipment);
  };

  render() {
    // TODO: replace minimal styling with actual styling during UI phase
    const { pageKey, pageList, match, push, newDutyStationAddress } = this.props;
    const { hasDeliveryAddress, initialValues, useCurrentResidence } = this.state;
    const fieldsetClasses = 'margin-top-2';
    return (
      <Formik initialValues={initialValues} enableReinitialize>
        {({ handleChange, values }) => (
          <WizardPage
            match={match}
            pageKey={pageKey}
            pageList={pageList}
            push={push}
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
                renderExistingAddressCheckbox={() => (
                  <Checkbox
                    data-testid="useCurrentResidence"
                    label="Use my current residence address"
                    name="useCurrentResidence"
                    checked={useCurrentResidence}
                    onChange={() => this.handleUseCurrentResidenceChange(values)}
                  />
                )}
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
                    <div>
                      {newDutyStationAddress.city}, {newDutyStationAddress.state} {newDutyStationAddress.postal_code}
                    </div>
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
                  data-testid="remarks"
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
  currentResidence: shape({
    street_address_1: string,
    street_address_2: string,
    state: string,
    postal_code: string,
  }).isRequired,
  pageKey: string.isRequired,
  pageList: arrayOf(string).isRequired,
  match: shape({
    isExact: bool.isRequired,
    params: shape({
      moveId: string.isRequired,
    }),
    path: string.isRequired,
    url: string.isRequired,
  }).isRequired,
  newDutyStationAddress: shape({
    city: string.isRequired,
    state: string.isRequired,
    post_code: string.isRequired,
  }).isRequired,
  moveTaskOrderID: string.isRequired,
  createMTOShipment: func.isRequired,
  showLoggedInUser: func.isRequired,
  push: func.isRequired,
};

const mapStateToProps = (state) => ({
  moveTaskOrderID: get(selectLoggedInUser(state), 'service_member.orders[0].move_task_order_id', ''),
  currentResidence: get(selectLoggedInUser(state), 'service_member.residential_address', {}),
  newDutyStationAddress: get(selectLoggedInUser(state), 'service_member.orders[0].new_duty_station.address', {}),
});

const mapDispatchToProps = {
  createMTOShipment: createMTOShipmentAction,
  showLoggedInUser: showLoggedInUserAction,
};

export { HHGDetailsForm as HHGDetailsFormComponent };
export default connect(mapStateToProps, mapDispatchToProps)(HHGDetailsForm);
