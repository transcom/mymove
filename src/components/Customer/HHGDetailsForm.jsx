import React, { Component } from 'react';
import { arrayOf, string, bool, shape, func } from 'prop-types';
import { get, isEmpty } from 'lodash';
import { connect } from 'react-redux';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
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
import { validateDate } from 'utils/formikValidators';

const PickupAddressSchema = Yup.object().shape({
  mailingAddress1: Yup.string().required('Required'),
  mailingAddress2: Yup.string(),
  city: Yup.string().required('Required'),
  state: Yup.string().length(2, 'Must use state abbreviation').required('Required'),
  zip: Yup.string()
    // eslint-disable-next-line security/detect-unsafe-regex
    .matches(/^(\d{5}([-]\d{4})?)$/, 'Must be valid zip code')
    .required('Required'),
});

const DeliveryAddressSchema = Yup.object().shape({
  mailingAddress1: Yup.string(),
  mailingAddress2: Yup.string(),
  city: Yup.string(),
  state: Yup.string().length(2, 'Must use state abbreviation'),
  zip: Yup.string()
    // eslint-disable-next-line security/detect-unsafe-regex
    .matches(/^(\d{5}([-]\d{4})?)$/, 'Must be valid zip code'),
});

const AgentSchema = Yup.object().shape({
  firstName: Yup.string(),
  lastName: Yup.string(),
  phone: Yup.string().matches(/^[2-9]\d{2}\d{3}\d{4}$/, 'Must be valid phone number'),
  email: Yup.string().email('Must be valid email'),
});
const HHGDetailsFormSchema = Yup.object().shape({
  // requiredPickupDate, requiredDeliveryDate are also required, but using field level validation
  pickupLocation: PickupAddressSchema,
  deliveryLocation: DeliveryAddressSchema,
  releasingAgent: AgentSchema,
  receivingAgent: AgentSchema,
  remarks: Yup.string(),
});
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
        state: pickupLocation.state.toUpperCase(),
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
        state: deliveryLocation.state.toUpperCase(),
        postal_code: deliveryLocation.zip,
        country: deliveryLocation.country,
      };
    }

    function formatAgent(agent) {
      const agentCopy = { ...agent };
      Object.keys(agentCopy).forEach((key) => {
        /* eslint-disable security/detect-object-injection */
        if (agentCopy[key] === '') {
          delete agentCopy[key];
        } else if (key === 'phone') {
          const phoneNum = agentCopy[key];
          // will be in format xxx-xxx-xxxx
          agentCopy[key] = `${phoneNum.slice(0, 3)}-${phoneNum.slice(3, 6)}-${phoneNum.slice(6, 10)}`;
        }
        /* eslint-enable security/detect-object-injection */
      });
      return agentCopy;
    }

    if (releasingAgent) {
      const formattedAgent = formatAgent(releasingAgent);
      if (!isEmpty(formattedAgent)) {
        mtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RELEASING });
      }
    }

    if (receivingAgent) {
      const formattedAgent = formatAgent(receivingAgent);
      if (!isEmpty(formattedAgent)) {
        mtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RECEIVING });
      }
    }
    createMTOShipment(mtoShipment);
  };

  render() {
    // TODO: replace minimal styling with actual styling during UI phase
    const { pageKey, pageList, match, push, newDutyStationAddress } = this.props;
    const { hasDeliveryAddress, initialValues, useCurrentResidence } = this.state;
    const fieldsetClasses = 'margin-top-2';
    return (
      <Formik
        initialValues={initialValues}
        enableReinitialize
        validateOnBlur
        validateOnChange
        validationSchema={HHGDetailsFormSchema}
      >
        {({ values, dirty, isValid }) => (
          <WizardPage
            canMoveNext={dirty && isValid}
            match={match}
            pageKey={pageKey}
            pageList={pageList}
            push={push}
            handleSubmit={() => this.submitMTOShipment(values)}
          >
            <Form>
              <Fieldset legend="Pickup date" className={fieldsetClasses}>
                <Field
                  as={DatePickerInput}
                  name="requestedPickupDate"
                  label="Requested pickup date"
                  id="requestedPickupDate"
                  value={values.requestedPickupDate}
                  validate={validateDate}
                />
              </Fieldset>
              <span className="usa-hint" id="pickupDateHint">
                Your movers will confirm this date or one shortly before or after.
              </span>

              <AddressFields
                name="pickupLocation"
                legend="Pickup location"
                className={fieldsetClasses}
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
                values={values.releasingAgent}
              />
              <Fieldset legend="Delivery date" className={fieldsetClasses}>
                <DatePickerInput
                  name="requestedDeliveryDate"
                  label="Requested delivery date"
                  id="requestedDeliveryDate"
                  value={values.requestedDeliveryDate}
                  validate={validateDate}
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
                  <AddressFields name="deliveryLocation" className={fieldsetClasses} values={values.deliveryLocation} />
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
                values={values.receivingAgent}
              />
              <Fieldset legend="Remarks" className={fieldsetClasses}>
                <TextInput
                  label="Anything else you would like us to know?"
                  labelHint="(optional)"
                  data-testid="remarks"
                  name="remarks"
                  id="remarks"
                  maxLength={1500}
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
    city: string,
    state: string,
    postal_code: string,
  }),
  moveTaskOrderID: string.isRequired,
  createMTOShipment: func.isRequired,
  showLoggedInUser: func.isRequired,
  push: func.isRequired,
};

HHGDetailsForm.defaultProps = {
  newDutyStationAddress: {
    city: '',
    state: '',
    postal_code: '',
  },
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
