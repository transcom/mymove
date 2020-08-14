import React, { Component } from 'react';
import { string, shape, func } from 'prop-types';
import { isEmpty } from 'lodash';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Fieldset, Radio, Label, Textarea, Button, Checkbox } from '@trussworks/react-uswds';

import { Form } from '../form/Form';
import { DatePickerInput } from '../form/fields';
import { AddressFields } from '../form/AddressFields/AddressFields';
import { ContactInfoFields } from '../form/ContactInfoFields/ContactInfoFields';

import { MTOAgentType } from 'shared/constants';
import { formatSwaggerDate } from 'shared/formatters';
import { validateDate } from 'utils/formikValidators';

import './EditShipment.scss';

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

const Divider = ({ className }) => <hr className={className} />;

Divider.propTypes = {
  className: string,
};

Divider.defaultProps = {
  className: '',
};

const Hint = ({ className, children, ...props }) => (
  // eslint-disable-next-line react/jsx-props-no-spreading
  <div {...props} className={`usa-hint ${className}`}>
    {children}
  </div>
);

Hint.propTypes = {
  className: string,
  children: React.Children.isRequired,
};

Hint.defaultProps = {
  className: '',
};

class EditShipment extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hasDeliveryAddress: true,
      initialValues: {},
    };
  }

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
    const { hasDeliveryAddress, initialValues, useCurrentResidence } = this.state;
    return (
      <div className="grid-container">
        <div className="margin-top-2 hhg-label">HHG</div>
        <h2 className="margin-top-1" style={{ fontSize: 28 }}>
          When and where will you move this shipment?
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
              <Fieldset legend="Pickup date" className="fieldset-legend margin-top-4">
                <Field
                  as={DatePickerInput}
                  name="requestedPickupDate"
                  label="Requested pickup date"
                  labelClassName="margin-top-2"
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
                className="margin-bottom-3 fieldset-legend"
                name="pickupLocation"
                legend="Pickup location"
                renderExistingAddressCheckbox={() => (
                  <Checkbox
                    className="margin-top-3"
                    data-testid="useCurrentResidence"
                    label="Use my current residence address"
                    name="useCurrentResidence"
                    checked={useCurrentResidence}
                    onChange={() => this.handleUseCurrentResidenceChange(values)}
                  />
                )}
                values={values.pickupLocation}
              />
              <div className="usa-hint">
                If you have more things at another pickup location, you can schedule for them later.
              </div>
              <hr className="margin-top-4 margin-bottom-4" />
              <ContactInfoFields
                className="fieldset-legend margin-bottom-5"
                name="releasingAgent"
                legend="Releasing agent"
                subtitle="Who can allow the movers to take your stuff if you're not there?"
                values={values.releasingAgent}
              />
              <Divider className="margin-bottom-6" />
              <Fieldset legend="Delivery date" className="fieldset-legend">
                <DatePickerInput
                  name="requestedDeliveryDate"
                  label="Requested delivery date"
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
              <Fieldset legend="Delivery location" className="fieldset-legend">
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
                className="fieldset-legend"
                name="receivingAgent"
                legend="Receiving agent"
                subtitle="Who can take delivery for you if the movers arrive and you're not there?"
                values={values.receivingAgent}
              />
              <Divider className="margin-top-4 margin-bottom-4" />
              <Fieldset legend="Remarks" className="fieldset-legend">
                <div className="small-bold margin-top-3 margin-bottom-1">
                  Is there anything special about this shipment that the movers should know?
                </div>
                <div className="hhg-examples-container">
                  <strong>Examples</strong>
                  <ul>
                    <li>Things that might need special handling</li>
                    <li>Access info for a location</li>
                    <li>Weapons or alcohol</li>
                  </ul>
                </div>

                <Textarea
                  label="Anything else you would like us to know?"
                  labelHint="(optional)"
                  data-testid="remarks"
                  name="remarks"
                  id="remarks"
                  placeholder="This is 500 characters of customer remarks placeholder"
                  maxLength={500}
                  style={{ resize: 'vertical' }}
                  type="textarea"
                  value={values.remarks}
                />
              </Fieldset>
              <Divider className="margin-top-6 margin-bottom-3" />
              <Hint className="margin-bottom-2">
                You can change details for your HHG shipment when you talk to your move counselor or the person
                who&apos;s your point of contact with the movers. You can also edit in MilMove up to 24 hours before
                your final pickup date.
              </Hint>
              <div style={{ display: 'flex', flexDirection: 'column' }}>
                <Button>Save</Button>
                <Button className="usa-button--unstyled">
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
  moveTaskOrderID: string.isRequired,
  createMTOShipment: func.isRequired,
};

export default EditShipment;
