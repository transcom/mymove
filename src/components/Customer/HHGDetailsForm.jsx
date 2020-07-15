import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import { Fieldset, Radio } from '@trussworks/react-uswds';

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
    // not currently getting hit
    this.setState((prevState) => {
      return { hasDeliveryAddress: !prevState.hasDeliveryAddress };
    });
  };

  render() {
    const { initialValues, pageKey, pageList } = this.props;
    const { hasDeliveryAddress } = this.state;
    return (
      <Formik initialValues={initialValues}>
        <WizardPage pageKey={pageKey} pageList={pageList} handleSubmit={() => {}}>
          <Form>
            <DatePickerInput name="requestedPickupDate" label="Requested pickup date" id="requested-pickup-date" />
            <AddressFields initialValues={initialValues.pickupLocation} legend="Pickup location" />
            <ContactInfoFields initialValues={initialValues.releasingAgent} legend="Releasing agent" />
            <DatePickerInput
              name="requestedDeliveryDate"
              label="Requested delivery date"
              id="requested-delivery-date"
            />
            <Fieldset legend="Delivery location">
              <Radio label="Yes" checked={hasDeliveryAddress} onChange={this.handleChangeHasDeliveryAddress} />
              <Radio label="No" checked={!hasDeliveryAddress} onChange={this.handleChangeHasDeliveryAddress} />
              {hasDeliveryAddress ? (
                <AddressFields initialValues={initialValues.deliveryLocation} />
              ) : (
                <>
                  <div>We can use the zip of your new duty station.</div>
                  <div>[City], [State] [New duty station zip]</div>
                </>
              )}
            </Fieldset>
            <ContactInfoFields initialValues={initialValues.receivingAgent} legend="Receiving agent" />
            <Fieldset legend="Remarks">
              <TextInput name="remarks" label="Remarks" id="requested-delivery-date" />
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
