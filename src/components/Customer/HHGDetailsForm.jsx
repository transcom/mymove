import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import { Fieldset, Radio, Label } from '@trussworks/react-uswds';

import { Form } from '../form/Form';
import { DatePickerInput, TextInput } from '../form/fields';
import { AddressFields } from '../form/AddressFields/AddressFields';
import { ContactInfoFields } from '../form/ContactInfoFields/ContactInfoFields';
import { WizardPage } from '../../shared/WizardPage';

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

  render() {
    // TODO: replace minimal styling with actual styling during UI phase
    const { initialValues, pageKey, pageList } = this.props;
    const { hasDeliveryAddress } = this.state;
    const fieldsetClasses = 'margin-top-2';
    return (
      <Formik initialValues={initialValues}>
        {({ handleChange, values }) => (
          <WizardPage pageKey={pageKey} pageList={pageList} handleSubmit={() => {}}>
            <Form>
              <Fieldset legend="Pickup date" className={fieldsetClasses}>
                <DatePickerInput
                  name="requestedPickupDate"
                  label="Requested pickup date"
                  id="requestedPickupDate"
                  onChange={handleChange}
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
                  handleChange={handleChange}
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
  pageKey: PropTypes.string.isRequired,
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
  initialValues: PropTypes.shape({
    requestedPickupDate: PropTypes.string,
    pickupLocation: PropTypes.shape({
      mailingAddress1: PropTypes.string,
      mailingAddress2: PropTypes.string,
      city: PropTypes.string,
      state: PropTypes.string,
      zip: PropTypes.string,
    }),
    requestedDeliveryDate: PropTypes.string,
    deliveryLocation: PropTypes.shape({
      mailingAddress1: PropTypes.string,
      mailingAddress2: PropTypes.string,
      city: PropTypes.string,
      state: PropTypes.string,
      zip: PropTypes.string,
    }),
    releasingAgent: PropTypes.shape({
      firstName: PropTypes.string,
      lastName: PropTypes.string,
      phone: PropTypes.string,
      email: PropTypes.string,
    }),
    receivingAgent: PropTypes.shape({
      firstName: PropTypes.string,
      lastName: PropTypes.string,
      phone: PropTypes.string,
      email: PropTypes.string,
    }),
    remarks: PropTypes.string,
  }),
};

HHGDetailsForm.defaultProps = {
  initialValues: {},
};

export default HHGDetailsForm;
