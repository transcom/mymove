import React, { Component } from 'react';
import { string, shape, func, bool } from 'prop-types';
import { isEmpty } from 'lodash';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Radio, Label, Textarea, Button, Checkbox } from '@trussworks/react-uswds';

import { Form } from '../form/Form';
import { DatePickerInput } from '../form/fields';
import { AddressFields } from '../form/AddressFields/AddressFields';
import { ContactInfoFields } from '../form/ContactInfoFields/ContactInfoFields';

import styles from './EditShipment.module.scss';

import { MTOAgentType, SHIPMENT_OPTIONS } from 'shared/constants';
import { formatSwaggerDate } from 'shared/formatters';
import { validateDate } from 'utils/formikValidators';
import Hint from 'shared/Hint';
import Fieldset from 'shared/Fieldset';
import Divider from 'shared/Divider';

const PickupAddressSchema = Yup.object().shape({
  street_address_1: Yup.string().required('Required'),
  street_address_2: Yup.string(),
  city: Yup.string().required('Required'),
  state: Yup.string().length(2, 'Must use state abbreviation').required('Required'),
  postal_code: Yup.string()
    // eslint-disable-next-line security/detect-unsafe-regex
    .matches(/^(\d{5}([-]\d{4})?)$/, 'Must be valid zip code')
    .required('Required'),
});

const DeliveryAddressSchema = Yup.object().shape({
  street_address_1: Yup.string(),
  street_address_2: Yup.string(),
  city: Yup.string(),
  state: Yup.string().length(2, 'Must use state abbreviation'),
  postal_code: Yup.string()
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
  pickupAddress: PickupAddressSchema,
  destinationAddress: DeliveryAddressSchema,
  releasingAgent: AgentSchema,
  receivingAgent: AgentSchema,
  customerRemarks: Yup.string(),
});

class EditShipment extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hasDeliveryAddress: false,
      initialValues: {},
    };
  }

  // componentDidMount() {
  //   const { showLoggedInUser, mtoShipment } = this.props;
  //   showLoggedInUser();
  //   if (mtoShipment.id) {
  //     this.setInitialState(mtoShipment);
  //   }
  // }

  // componentDidUpdate(prevProps) {
  //   const { mtoShipment } = this.props;
  //   // If refreshing page, need to handle mtoShipment populating from a promise
  //   if (mtoShipment.id && prevProps.mtoShipment.id !== mtoShipment.id) {
  //     this.setInitialState(mtoShipment);
  //   }
  // }

  // setInitialState = (mtoShipment) => {
  //   function cleanAgentPhone(agent) {
  //     const agentCopy = { ...agent };
  //     Object.keys(agentCopy).forEach((key) => {
  //       /* eslint-disable security/detect-object-injection */
  //       if (key === 'phone') {
  //         const phoneNum = agentCopy[key];
  //         // will be in format xxxxxxxxxx
  //         agentCopy[key] = phoneNum.split('-').join('');
  //       }
  //     });
  //     return agentCopy;
  //   }
  //   // for existing mtoShipment, reshape agents from array of objects to key/object for proper handling
  //   const { agents } = mtoShipment;
  //   const formattedMTOShipment = { ...mtoShipment };
  //   if (agents) {
  //     const receivingAgent = agents.find((agent) => agent.agentType === 'RECEIVING_AGENT');
  //     const releasingAgent = agents.find((agent) => agent.agentType === 'RELEASING_AGENT');

  //     // Remove dashes from agent phones for expected form phone format
  //     if (receivingAgent) {
  //       const formattedAgent = cleanAgentPhone(releasingAgent);
  //       if (!isEmpty(formattedAgent)) {
  //         formattedMTOShipment.receivingAgent = { ...formattedAgent };
  //       }
  //     }
  //     if (releasingAgent) {
  //       const formattedAgent = cleanAgentPhone(releasingAgent);
  //       if (!isEmpty(formattedAgent)) {
  //         formattedMTOShipment.releasingAgent = { ...formattedAgent };
  //       }
  //     }
  //   }
  //   this.setState({ initialValues: formattedMTOShipment });
  // };

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
              pickupAddress: {
                street_address_1: currentResidence.street_address_1,
                street_address_2: currentResidence.street_address_2,
                city: currentResidence.city,
                state: currentResidence.state,
                postal_code: currentResidence.postal_code,
              },
            },
          });
        } else {
          this.setState({
            initialValues: {
              ...initialValues,
              ...currentValues,
              pickupAddress: {
                street_address_1: '',
                street_address_2: '',
                city: '',
                state: '',
                postal_code: '',
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
    pickupAddress,
    destinationAddress,
    receivingAgent,
    releasingAgent,
    customerRemarks,
  }) => {
    const { createMTOShipment, match } = this.props;
    const { hasDeliveryAddress } = this.state;
    const { moveId } = match.params;
    const pendingMtoShipment = {
      moveTaskOrderID: moveId,
      shipmentType: SHIPMENT_OPTIONS.HHG,
      requestedPickupDate: formatSwaggerDate(requestedPickupDate),
      requestedDeliveryDate: formatSwaggerDate(requestedDeliveryDate),
      customerRemarks,
      pickupAddress: {
        street_address_1: pickupAddress.street_address_1,
        street_address_2: pickupAddress.street_address_2,
        city: pickupAddress.city,
        state: pickupAddress.state.toUpperCase(),
        postal_code: pickupAddress.postal_code,
        country: pickupAddress.country,
      },
      agents: [],
    };

    if (hasDeliveryAddress) {
      pendingMtoShipment.destinationAddress = {
        street_address_1: destinationAddress.street_address_1,
        street_address_2: destinationAddress.street_address_2,
        city: destinationAddress.city,
        state: destinationAddress.state.toUpperCase(),
        postal_code: destinationAddress.postal_code,
        country: destinationAddress.country,
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
        pendingMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RELEASING });
      }
    }

    if (receivingAgent) {
      const formattedAgent = formatAgent(receivingAgent);
      if (!isEmpty(formattedAgent)) {
        pendingMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RECEIVING });
      }
    }
    createMTOShipment(pendingMtoShipment);
  };

  render() {
    // const { pageKey, pageList, match, push, newDutyStationAddress, mtoShipment } = this.props;

    const { hasDeliveryAddress, initialValues, useCurrentResidence } = this.state;
    return (
      <div className="grid-container">
        <div className={`margin-top-2 ${styles['hhg-label']}`}>HHG</div>
        <h2 className="margin-top-1" style={{ fontSize: 28 }}>
          When and where can the movers pick up and deliver this shipment?
        </h2>
        <Formik
          initialValues={initialValues}
          enableReinitialize
          validateOnBlur
          validateOnChange
          validationSchema={HHGDetailsFormSchema}
        >
          {({ values }) => (
            <Form>
              <Fieldset legend="Pickup date" className="margin-top-4">
                <Field
                  as={DatePickerInput}
                  name="requestedPickupDate"
                  label="Requested pickup date"
                  labelClassName={`margin-top-2 ${styles['small-bold']}`}
                  id="requestedPickupDate"
                  value={values.requestedPickupDate}
                  validate={validateDate}
                />
              </Fieldset>
              <Hint className="margin-top-1" id="pickupDateHint">
                Movers will contact you to schedule the actual pickup date. That date should fall within 7 days of your
                requested date. Tip: Avoid scheduling multiple shipments on the same day.{' '}
              </Hint>
              <Divider className="margin-top-4 margin-bottom-4" />
              <AddressFields
                className="margin-bottom-3"
                name="pickupAddress"
                legend="Pickup location"
                renderExistingAddressCheckbox={() => (
                  <Checkbox
                    className="margin-top-3"
                    data-testid="useCurrentResidence"
                    label="Use my current address"
                    name="useCurrentResidence"
                    checked={useCurrentResidence}
                    onChange={() => this.handleUseCurrentResidenceChange(values)}
                  />
                )}
                values={values.pickupAddress}
              />
              <Hint>If you have more things at another pickup location, you can schedule for them later.</Hint>
              <hr className="margin-top-4 margin-bottom-4" />
              <ContactInfoFields
                className="margin-bottom-5"
                name="releasingAgent"
                legend="Releasing agent"
                hintText="Optional"
                subtitle="Who can allow the movers to take your stuff if you're not there?"
                subtitleClassName="margin-top-3"
                values={values.releasingAgent}
              />
              <Divider className="margin-bottom-6" />
              <Fieldset legend="Delivery date">
                <DatePickerInput
                  name="requestedDeliveryDate"
                  label="Requested delivery date"
                  labelClassName={`${styles['small-bold']}`}
                  id="requestedDeliveryDate"
                  value={values.requestedDeliveryDate}
                  validate={validateDate}
                />
                <Hint className="margin-top-1">
                  Shipments can take several weeks to arrive, depending on how far they&apos;re going. Your movers will
                  contact you close to the date you select to coordinate delivery
                </Hint>
              </Fieldset>
              <Divider className="margin-top-4 margin-bottom-4" />
              <Fieldset legend="Delivery location">
                <Label className="margin-top-3 margin-bottom-1">Do you know your delivery address?</Label>
                <div className="display-flex margin-top-1">
                  <Radio
                    className="margin-right-3"
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
                {hasDeliveryAddress && <AddressFields name="deliveryLocation" values={values.deliveryLocation} />}
              </Fieldset>
              <Divider className="margin-top-4 margin-bottom-4" />
              <ContactInfoFields
                name="receivingAgent"
                legend="Receiving agent"
                hintText="Optional"
                subtitleClassName="margin-top-3"
                subtitle="Who can take delivery for you if the movers arrive and you're not there?"
                values={values.receivingAgent}
              />
              <Divider className="margin-top-4 margin-bottom-4" />
              <Fieldset hintText="Optional" legend="Remarks">
                <div className={`${styles['small-bold']} margin-top-3 margin-bottom-1`}>
                  Is there anything special about this shipment that the movers should know?
                </div>
                <div className={`${styles['hhg-examples-container']}`}>
                  <strong>Examples</strong>
                  <ul>
                    <li>Things that might need special handling</li>
                    <li>Access info for a location</li>
                    <li>Weapons or alcohol</li>
                  </ul>
                </div>

                <Textarea
                  label="Anything else you would like us to know?"
                  data-testid="remarks"
                  name="customerRemarks"
                  className={`${styles.remarks}`}
                  placeholder="This is 500 characters of customer remarks placeholder"
                  maxLength={500}
                  type="textarea"
                  value={values.customerRemarks}
                />
              </Fieldset>
              <Divider className="margin-top-6 margin-bottom-3" />
              <Hint className="margin-bottom-2">
                You can change details for your HHG shipment when you talk to your move counselor or the person
                who&apos;s your point of contact with the movers. You can also edit in MilMove up to 24 hours before
                your final pickup date.
              </Hint>
              <div style={{ display: 'flex', flexDirection: 'column' }}>
                <Button>
                  <span>Save</span>
                </Button>
                <Button className={`${styles['cancel-button']}`}>
                  <span>Cancel</span>
                </Button>
              </div>
            </Form>
          )}
        </Formik>
      </div>
    );
  }
}

EditShipment.propTypes = {
  currentResidence: shape({
    street_address_1: string,
    street_address_2: string,
    state: string,
    postal_code: string,
  }).isRequired,
  // moveTaskOrderID: string.isRequired,
  createMTOShipment: func.isRequired,
  match: shape({
    isExact: bool.isRequired,
    params: shape({
      moveId: string.isRequired,
    }),
    path: string.isRequired,
    url: string.isRequired,
  }).isRequired,
};

export default EditShipment;
