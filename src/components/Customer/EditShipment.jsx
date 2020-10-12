import React, { Component } from 'react';
import { arrayOf, string, shape, func, bool } from 'prop-types';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Radio, Label, Textarea, Button } from '@trussworks/react-uswds';

import { Form } from '../form/Form';
import { DatePickerInput } from '../form/fields';
import { AddressFields } from '../form/AddressFields/AddressFields';
import { ContactInfoFields } from '../form/ContactInfoFields/ContactInfoFields';
import { SectionWrapper } from 'Containers/SectionWrapper';

import styles from './EditShipment.module.scss';

import {
  selectMTOShipmentById,
  updateMTOShipment as updateMTOShipmentAction,
} from 'shared/Entities/modules/mtoShipments';
import Checkbox from 'shared/Checkbox';
import { selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';
import { MTOAgentType, SHIPMENT_OPTIONS } from 'shared/constants';
import { formatSwaggerDate } from 'shared/formatters';
import { validateDate } from 'utils/formikValidators';
import Hint from 'shared/Hint';
import Fieldset from 'shared/Fieldset';

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
      useCurrentResidence: false,
      initialValues: {},
    };
  }

  componentDidMount() {
    const { mtoShipment } = this.props;
    if (mtoShipment.id) {
      this.setInitialState(mtoShipment);
    }
  }

  componentDidUpdate(prevProps) {
    const { mtoShipment } = this.props;
    // If refreshing page, need to handle mtoShipment populating from a promise
    if (mtoShipment.id && prevProps.mtoShipment.id !== mtoShipment.id) {
      this.setInitialState(mtoShipment);
    }
  }

  setInitialState = (mtoShipment) => {
    function cleanAgentPhone(agent) {
      const agentCopy = { ...agent };
      Object.keys(agentCopy).forEach((key) => {
        /* eslint-disable security/detect-object-injection */
        if (key === 'phone') {
          const phoneNum = agentCopy[key];
          // will be in format xxxxxxxxxx
          agentCopy[key] = phoneNum.split('-').join('');
        }
      });
      return agentCopy;
    }
    // for existing mtoShipment, reshape agents from array of objects to key/object for proper handling
    const { agents } = mtoShipment;
    const formattedMTOShipment = { ...mtoShipment };
    if (agents) {
      const receivingAgent = agents.find((agent) => agent.agentType === 'RECEIVING_AGENT');
      const releasingAgent = agents.find((agent) => agent.agentType === 'RELEASING_AGENT');

      // Remove dashes from agent phones for expected form phone format
      if (receivingAgent) {
        const formattedAgent = cleanAgentPhone(receivingAgent);
        if (Object.keys(formattedAgent).length) {
          formattedMTOShipment.receivingAgent = { ...formattedAgent };
        }
      }
      if (releasingAgent) {
        const formattedAgent = cleanAgentPhone(releasingAgent);
        if (Object.keys(formattedAgent).length) {
          formattedMTOShipment.releasingAgent = { ...formattedAgent };
        }
      }
    }
    const hasDeliveryAddress = get(mtoShipment, 'destinationAddress', false);
    this.setState({ initialValues: formattedMTOShipment, hasDeliveryAddress });
  };

  handleChangeHasDeliveryAddress = () => {
    this.setState((prevState) => {
      return { hasDeliveryAddress: !prevState.hasDeliveryAddress };
    });
  };

  // Use current residence
  handleUseCurrentResidenceChange = (currentValues) => {
    // eslint-disable-next-line react/destructuring-assignment
    const { initialValues } = this.state;
    const { currentResidence, match, mtoShipment } = this.props;
    this.setState(
      (state) => ({ useCurrentResidence: !state.useCurrentResidence }),
      () => {
        // eslint-disable-next-line react/destructuring-assignment
        if (this.state.useCurrentResidence) {
          this.setState({
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
          // eslint-disable-next-line no-lonely-if
          if (match.params.moveId === initialValues.moveTaskOrderID) {
            this.setState({
              initialValues: {
                ...initialValues,
                ...currentValues,
                pickupAddress: {
                  street_address_1: mtoShipment.pickupAddress.street_address_1,
                  street_address_2: mtoShipment.pickupAddress.street_address_2,
                  city: mtoShipment.pickupAddress.city,
                  state: mtoShipment.pickupAddress.state,
                  postal_code: mtoShipment.pickupAddress.postal_code,
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
    const { updateMTOShipment, match, mtoShipment, history } = this.props;
    const { hasDeliveryAddress } = this.state;
    const { moveId } = match.params;
    const goBack = get(history, 'goBack', '');
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
      if (Object.keys(formattedAgent).length) {
        pendingMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RELEASING });
      }
    }

    if (receivingAgent) {
      const formattedAgent = formatAgent(receivingAgent);
      if (Object.keys(formattedAgent).length) {
        pendingMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RECEIVING });
      }
    }
    updateMTOShipment(mtoShipment.id, pendingMtoShipment, mtoShipment.eTag).then(() => {
      goBack();
    });
  };

  getShipmentNumber = () => {
    const { search } = window.location;
    const params = new URLSearchParams(search);
    const shipmentNumber = params.get('shipmentNumber');
    return shipmentNumber;
  };

  render() {
    const { history } = this.props;
    const { hasDeliveryAddress, initialValues, useCurrentResidence } = this.state;
    const goBack = get(history, 'goBack', '');
    const shipmentNumber = this.getShipmentNumber();

    return (
      <div className="grid-container">
        <div className={`margin-top-2 ${styles['hhg-label']}`}>{`HHG ${!!shipmentNumber && shipmentNumber}`}</div>
        <h1 className="margin-top-1">When and where can the movers pick up and deliver this shipment?</h1>
        <Formik
          initialValues={initialValues}
          enableReinitialize
          validateOnBlur
          validateOnChange
          validationSchema={HHGDetailsFormSchema}
        >
          {({ values, isValid, dirty, isSubmitting, handleChange }) => (
            <div className="wrapper-co">
              <Form>
                <SectionWrapper>
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
                    <Hint className="margin-top-1" id="pickupDateHint">
                      Movers will contact you to schedule the actual pickup date. That date should fall within 7 days of
                      your requested date. Tip: Avoid scheduling multiple shipments on the same day.{' '}
                    </Hint>
                  </Fieldset>
                </SectionWrapper>
                <SectionWrapper>
                  <AddressFields
                    className="margin-bottom-3"
                    name="pickupAddress"
                    legend="Pickup location"
                    renderExistingAddressCheckbox={() => (
                      <div className="margin-y-2">
                        <Checkbox
                          data-testid="useCurrentResidence"
                          label="Use my current residence address"
                          name="useCurrentResidence"
                          checked={useCurrentResidence}
                          onChange={() => this.handleUseCurrentResidenceChange(values)}
                        />
                      </div>
                    )}
                    values={values.pickupAddress}
                  />
                  <Hint>If you have more things at another pickup location, you can schedule for them later.</Hint>
                </SectionWrapper>
                <SectionWrapper>
                  <ContactInfoFields
                    className="margin-bottom-5"
                    name="releasingAgent"
                    legend="Releasing agent"
                    hintText="Optional"
                    subtitle="Who can allow the movers to take your stuff if you're not there?"
                    subtitleClassName="margin-top-3"
                    values={values.releasingAgent}
                  />
                </SectionWrapper>
                <SectionWrapper>
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
                      Shipments can take several weeks to arrive, depending on how far they&apos;re going. Your movers
                      will contact you close to the date you select to coordinate delivery
                    </Hint>
                  </Fieldset>
                </SectionWrapper>
                <SectionWrapper>
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
                    {hasDeliveryAddress && (
                      <AddressFields name="destinationAddress" values={values.destinationAddress} />
                    )}
                  </Fieldset>
                </SectionWrapper>
                <SectionWrapper>
                  <ContactInfoFields
                    name="receivingAgent"
                    legend="Receiving agent"
                    hintText="Optional"
                    subtitleClassName="margin-top-3"
                    subtitle="Who can take delivery for you if the movers arrive and you're not there?"
                    values={values.receivingAgent}
                  />
                </SectionWrapper>
                <SectionWrapper>
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
                      placeholder="500 characters"
                      maxLength={500}
                      type="textarea"
                      value={values.customerRemarks}
                      onChange={handleChange}
                    />
                  </Fieldset>
                </SectionWrapper>
                <Hint className="margin-bottom-2">
                  You can change details for your HHG shipment when you talk to your move counselor or the person
                  who&apos;s your point of contact with the movers. You can also edit in MilMove up to 24 hours before
                  your final pickup date.
                </Hint>
                <div style={{ display: 'flex', flexDirection: 'column' }}>
                  <Button
                    disabled={isSubmitting || (!isValid && !dirty) || (isValid && !dirty)}
                    onClick={() => this.submitMTOShipment(values)}
                  >
                    <span>Save</span>
                  </Button>
                  <Button className={`${styles['cancel-button']}`} onClick={goBack}>
                    <span>Cancel</span>
                  </Button>
                </div>
              </Form>
            </div>
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
  updateMTOShipment: func.isRequired,
  match: shape({
    isExact: bool.isRequired,
    params: shape({
      moveId: string.isRequired,
    }),
    path: string.isRequired,
    url: string.isRequired,
  }).isRequired,
  history: shape({
    goBack: func.isRequired,
  }).isRequired,
  mtoShipment: shape({
    agents: arrayOf(
      shape({
        firstName: string,
        lastName: string,
        phone: string,
        email: string,
        agentType: string,
      }),
    ),
    customerRemarks: string,
    requestedPickupDate: string,
    requestedDeliveryDate: string,
    pickupAddress: shape({
      city: string,
      postal_code: string,
      state: string,
      street_address_1: string,
    }),
    destinationAddress: shape({
      city: string,
      postal_code: string,
      state: string,
      street_address_1: string,
    }),
  }),
};

EditShipment.defaultProps = {
  mtoShipment: {
    id: '',
    customerRemarks: '',
    requestedPickupDate: '',
    requestedDeliveryDate: '',
    destinationAddress: {
      city: '',
      postal_code: '',
      state: '',
      street_address_1: '',
    },
  },
};

const mapStateToProps = (state, ownProps) => {
  const orders = selectActiveOrLatestOrdersFromEntities(state);

  const props = {
    mtoShipment: selectMTOShipmentById(state, ownProps.match.params.mtoShipmentId),
    currentResidence: get(selectServiceMemberFromLoggedInUser(state), 'residential_address', {}),
    newDutyStationAddress: get(orders, 'new_duty_station.address', {}),
  };
  return props;
};

const mapDispatchToProps = {
  updateMTOShipment: updateMTOShipmentAction,
};

export { EditShipment as EditShipmentComponent };
export default connect(mapStateToProps, mapDispatchToProps)(EditShipment);
