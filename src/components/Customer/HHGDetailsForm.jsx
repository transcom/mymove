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

  handleChangeHasDeliveryAddress = () => {
    this.setState((prevState) => {
      return { hasDeliveryAddress: !prevState.hasDeliveryAddress };
    });
  };

  render() {
    const { initialValues, pageKey, pageList } = this.props;
    const { hasDeliveryAddress } = this.state;
    const fieldsetClasses = { margin: 'margin-top-2' };
    return (
      <Formik initialValues={initialValues}>
        <WizardPage pageKey={pageKey} pageList={pageList} handleSubmit={() => {}}>
          <Form>
            <Fieldset legend="Pickup date" className={fieldsetClasses.margin}>
              <DatePickerInput name="requestedPickupDate" label="Requested pickup date" id="requested-pickup-date" />
            </Fieldset>
            <span className="usa-hint" id="pickupDateHint">
              Your movers will confirm this date or one shortly before or after.
            </span>
            <AddressFields
              initialValues={initialValues.pickupLocation}
              legend="Pickup location"
              className={fieldsetClasses.margin}
            />
            <ContactInfoFields
              initialValues={initialValues.releasingAgent}
              legend="Releasing agent"
              className={fieldsetClasses.margin}
              subtitle="Who can allow the movers to take your stuff if you're not there?"
            />
            <DatePickerInput
              name="requestedDeliveryDate"
              label="Requested delivery date"
              id="requested-delivery-date"
            />
            <span className="usa-hint" id="deliveryDateHint">
              Your movers will confirm this date or one shortly before or after.
            </span>
            <Fieldset legend="Delivery location" className={fieldsetClasses.margin}>
              <Label>Do you know your delivery address?</Label>
              <Radio
                className="display-inline"
                id="has-delivery-address"
                label="Yes"
                name="has-delivery-address"
                onChange={this.handleChangeHasDeliveryAddress}
              />
              <Radio
                className="display-inline-flex"
                id="no-delivery-address"
                label="No"
                name="has-delivery-address"
                defaultChecked
                onChange={this.handleChangeHasDeliveryAddress}
              />
              {hasDeliveryAddress ? (
                <AddressFields initialValues={initialValues.deliveryLocation} className={fieldsetClasses.margin} />
              ) : (
                <>
                  <div className={fieldsetClasses.margin}>We can use the zip of your new duty station.</div>
                  <div>[City], [State] [New duty station zip]</div>
                </>
              )}
            </Fieldset>
            <ContactInfoFields
              initialValues={initialValues.receivingAgent}
              legend="Receiving agent"
              className={fieldsetClasses.margin}
              subtitle="Who can take delivery for you if the movers arrive and you're not there?"
            />
            <Fieldset legend="Remarks" className={fieldsetClasses.margin}>
              <Label hint="(optional)">Anything else you would like us to know?</Label>
              <TextInput name="remarks" id="requested-delivery-date" maxLength={1500} />
            </Fieldset>
          </Form>
        </WizardPage>
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
