import React, { Component } from 'react';
import { bool, string, func } from 'prop-types';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Fieldset, Radio, Label } from '@trussworks/react-uswds';

import styles from './MtoShipmentForm.module.scss';
import { RequiredPlaceSchema, OptionalPlaceSchema } from './validationSchemas';

import { DatePickerInput, TextInput } from 'components/form/fields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { Form } from 'components/form/Form';
import { selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { WizardPage } from 'shared/WizardPage';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import Checkbox from 'shared/Checkbox';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { HhgShipmentShape, MtoDisplayOptionsShape } from 'types/customerShapes';
import { formatMtoShipment } from 'utils/formatMtoShipment';
import { validateDate } from 'utils/formikValidators';

class MtoShipmentInnerForm extends Component {
  render() {
    const {
      displayOptions,
      shipmentNumber,
      fieldsetClasses,
      values,
      hasDeliveryAddress,
      onHasDeliveryAddressChange,
      useCurrentResidence,
      onUseCurrentResidenceChange,
    } = this.props;

    return (
      <div>
        <div className={`margin-top-2 ${styles['hhg-label']}`}>
          {`${displayOptions.displayName} ${!!shipmentNumber ? shipmentNumber : ''}`}
        </div>
        <h1 className="margin-top-1">When and where can the movers pick up and deliver this shipment?</h1>
        <Form className={styles.HHGDetailsForm}>
          {displayOptions.showPickupFields && (
            <div>
              <Fieldset legend="Pickup date" className={fieldsetClasses}>
                <Field
                  as={DatePickerInput}
                  name="pickup.requestedDate"
                  label="Requested pickup date"
                  id="requestedPickupDate"
                  value={values.pickup.requestedDate}
                  validate={validateDate}
                />
                <span className="usa-hint" id="pickupDateHint">
                  Your movers will confirm this date or one shortly before or after.
                </span>
              </Fieldset>

              <AddressFields
                name="pickup.address"
                legend="Pickup location"
                className={fieldsetClasses}
                renderExistingAddressCheckbox={() => (
                  <div className="margin-y-2">
                    <Checkbox
                      data-testid="useCurrentResidence"
                      label="Use my current residence address"
                      name="useCurrentResidence"
                      checked={useCurrentResidence}
                      onChange={() => onUseCurrentResidenceChange(values)}
                    />
                  </div>
                )}
                values={values.pickup.address}
              />
              <ContactInfoFields
                name="pickup.agent"
                legend="Releasing agent"
                className={fieldsetClasses}
                subtitle="Who can allow the movers to take your stuff if you're not there?"
                subtitleClassName="margin-y-2"
                values={values.pickup.agent}
              />
            </div>
          )}
          {displayOptions.showDeliveryFields && (
            <div>
              <Fieldset legend="Delivery date" className={fieldsetClasses}>
                <DatePickerInput
                  name="delivery.requestedDate"
                  label="Requested delivery date"
                  id="requestedDeliveryDate"
                  value={values.delivery.requestedDate}
                  validate={validateDate}
                />
                <small className="usa-hint" id="deliveryDateHint">
                  Your movers will confirm this date or one shortly before or after.
                </small>
              </Fieldset>
              <Fieldset legend="Delivery location" className={fieldsetClasses}>
                <Label>Do you know your delivery address?</Label>
                <div className="display-flex margin-top-1">
                  <Radio
                    id="has-delivery-address"
                    label="Yes"
                    name="hasDeliveryAddress"
                    onChange={onHasDeliveryAddressChange}
                    checked={hasDeliveryAddress}
                  />
                  <Radio
                    id="no-delivery-address"
                    label="No"
                    name="hasDeliveryAddress"
                    checked={!hasDeliveryAddress}
                    onChange={onHasDeliveryAddressChange}
                  />
                </div>
                {hasDeliveryAddress ? (
                  <AddressFields name="destinationAddress" values={values.delivery.address} />
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
              <ContactInfoFields
                name="delivery.agent"
                legend="Receiving agent"
                className={fieldsetClasses}
                subtitle="Who can take delivery for you if the movers arrive and you're not there?"
                subtitleClassName="margin-y-2"
                values={values.delivery.agent}
              />
            </div>
          )}
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
        </Form>
      </div>
    );
  }
}

MtoShipmentInnerForm.propTypes = {
  currentResidence: AddressShape.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  selectedMoveType: string.isRequired,
  displayOptions: MtoDisplayOptionsShape.isRequired,
  mtoShipment: HhgShipmentShape,
  hasDeliveryAddress: bool.isRequired,
  onHasDeliveryAddressChange: func.isRequired,
  useCurrentResidence: bool.isRequired,
  onUseCurrentResidenceChange: func.isRequired,
};

MtoShipmentInnerForm.defaultProps = {
  wizardPage: {
    pageList: [],
    pageKey: '',
    match: { isExact: false, params: { moveID: '' } },
  },
  newDutyStationAddress: {
    city: '',
    state: '',
    postal_code: '',
  },
  deliveryOptions: {
    schema: {},
    displayName: '',
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
