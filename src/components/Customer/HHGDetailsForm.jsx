import React, { Component } from 'react';
import { arrayOf, string, bool, shape, func } from 'prop-types';
import { get, isEmpty } from 'lodash';
import { connect } from 'react-redux';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Fieldset, Radio, Label } from '@trussworks/react-uswds';

import { DatePickerInput, TextInput } from '../form/fields';
import { ContactInfoFields } from '../form/ContactInfoFields/ContactInfoFields';
import { AddressFields } from '../form/AddressFields/AddressFields';
import { Form } from '../form/Form';

import styles from './HHGDetailsForm.module.scss';
import SectionWrapper from './SectionWrapper';

import {
  selectMTOShipmentForMTO,
  createMTOShipment as createMTOShipmentAction,
  updateMTOShipment as updateMTOShipmentAction,
} from 'shared/Entities/modules/mtoShipments';
import { selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { WizardPage } from 'shared/WizardPage/index';
import { MTOAgentType, SHIPMENT_OPTIONS } from 'shared/constants';
import { formatSwaggerDate } from 'shared/formatters';
import Checkbox from 'shared/Checkbox';
import { validateDate } from 'utils/formikValidators';

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
class HHGDetailsForm extends Component {
  constructor(props) {
    super(props);
    const hasDeliveryAddress = get(props.mtoShipment, 'destinationAddress', false);
    this.state = {
      hasDeliveryAddress,
      useCurrentResidence: false,
      initialValues: {},
    };
  }

  componentDidMount() {
    const { showLoggedInUser } = this.props;
    showLoggedInUser();
  }

  // Use current residence
  handleUseCurrentResidenceChange = (currentValues) => {
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

  handleChangeHasDeliveryAddress = () => {
    this.setState((prevState) => {
      return { hasDeliveryAddress: !prevState.hasDeliveryAddress };
    });
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
    // TODO: replace minimal styling with actual styling during UI phase
    const { pageKey, pageList, match, push, newDutyStationAddress } = this.props;
    const { hasDeliveryAddress, useCurrentResidence, initialValues } = this.state;
    const fieldsetClasses = 'margin-top-2';
    return (
      <div className="wrapper-co">
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
              handleSubmit={() => this.submitMTOShipment(values, dirty)}
            >
              <h1>Now lets arrange details for the professional movers</h1>
              <Form className={styles.HHGDetailsForm}>
                <SectionWrapper>
                  <Fieldset legend="Pickup date" className={fieldsetClasses}>
                    <Field
                      as={DatePickerInput}
                      name="requestedPickupDate"
                      label="Requested pickup date"
                      id="requestedPickupDate"
                      value={values.requestedPickupDate}
                      validate={validateDate}
                    />
                    <span className="usa-hint" id="pickupDateHint">
                      Your movers will confirm this date or one shortly before or after.
                    </span>
                  </Fieldset>
                </SectionWrapper>

                <SectionWrapper>
                  <AddressFields
                    name="pickupAddress"
                    legend="Pickup location"
                    className={fieldsetClasses}
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
                </SectionWrapper>

                <SectionWrapper>
                  <ContactInfoFields
                    name="releasingAgent"
                    legend="Releasing agent"
                    className={fieldsetClasses}
                    subtitle="Who can allow the movers to take your stuff if you're not there?"
                    subtitleClassName="margin-y-2"
                    values={values.releasingAgent}
                  />
                </SectionWrapper>

                <SectionWrapper>
                  <Fieldset legend="Delivery date" className={fieldsetClasses}>
                    <DatePickerInput
                      name="requestedDeliveryDate"
                      label="Requested delivery date"
                      id="requestedDeliveryDate"
                      value={values.requestedDeliveryDate}
                      validate={validateDate}
                    />
                    <small className="usa-hint" id="deliveryDateHint">
                      Your movers will confirm this date or one shortly before or after.
                    </small>
                  </Fieldset>
                </SectionWrapper>

                <SectionWrapper>
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
                      <AddressFields name="destinationAddress" values={values.destinationAddress} />
                    ) : (
                      <>
                        <div>
                          <p className={fieldsetClasses}>
                            We can use the zip of your new duty station.
                            <br />
                            <strong>
                              {newDutyStationAddress.city}, {newDutyStationAddress.state}{' '}
                              {newDutyStationAddress.postal_code}{' '}
                            </strong>
                          </p>
                        </div>
                      </>
                    )}
                  </Fieldset>
                </SectionWrapper>

                <SectionWrapper>
                  <ContactInfoFields
                    name="receivingAgent"
                    legend="Receiving agent"
                    className={fieldsetClasses}
                    subtitle="Who can take delivery for you if the movers arrive and you're not there?"
                    subtitleClassName="margin-y-2"
                    values={values.receivingAgent}
                  />
                </SectionWrapper>

                <SectionWrapper>
                  <Fieldset legend="Remarks" className={fieldsetClasses}>
                    <TextInput
                      label="Anything else you would like us to know?"
                      labelHint="(optional)"
                      data-testid="remarks"
                      name="customerRemarks"
                      id="customerRemarks"
                      maxLength={1500}
                      value={values.customerRemarks}
                    />
                  </Fieldset>
                </SectionWrapper>
              </Form>
            </WizardPage>
          )}
        </Formik>
      </div>
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
  createMTOShipment: func.isRequired,
  showLoggedInUser: func.isRequired,
  push: func.isRequired,
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

HHGDetailsForm.defaultProps = {
  newDutyStationAddress: {
    city: '',
    state: '',
    postal_code: '',
  },
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
    mtoShipment: selectMTOShipmentForMTO(state, ownProps.match.params.moveId),
    currentResidence: get(selectServiceMemberFromLoggedInUser(state), 'residential_address', {}),
    newDutyStationAddress: get(orders, 'new_duty_station.address', {}),
  };
  return props;
};

const mapDispatchToProps = {
  createMTOShipment: createMTOShipmentAction,
  updateMTOShipment: updateMTOShipmentAction,
  showLoggedInUser: showLoggedInUserAction,
};

export { HHGDetailsForm as HHGDetailsFormComponent };
export default connect(mapStateToProps, mapDispatchToProps)(HHGDetailsForm);
