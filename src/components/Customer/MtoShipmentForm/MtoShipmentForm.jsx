import React, { Component } from 'react';
import { string, func } from 'prop-types';
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
import {
  selectMTOShipmentForMTO,
  createMTOShipment as createMTOShipmentAction,
  updateMTOShipment as updateMTOShipmentAction,
} from 'shared/Entities/modules/mtoShipments';
import { selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { WizardPage } from 'shared/WizardPage';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import Checkbox from 'shared/Checkbox';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { HhgShipmentShape, MatchShape, HistoryShape, PageKeyShape, PageListShape } from 'types/customerShapes';
import { formatMtoShipment } from 'utils/formatMtoShipment';
import { validateDate } from 'utils/formikValidators';

const hhgShipmentSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  delivery: OptionalPlaceSchema,
  customerRemarks: Yup.string(),
});

const ntsShipmentSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  customerRemarks: Yup.string(),
});

const ntsReleaseShipmentSchema = Yup.object().shape({
  delivery: OptionalPlaceSchema,
  customerRemarks: Yup.string(),
});

function getShipmentOptions(shipmentType) {
  switch (shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
      return {
        schema: hhgShipmentSchema,
        showPickupFields: true,
        showDeliveryFields: true,
      };
    case SHIPMENT_OPTIONS.NTS:
      return {
        schema: ntsShipmentSchema,
        showPickupFields: true,
        showDeliveryFields: false,
      };
    case SHIPMENT_OPTIONS.NTSR:
      return {
        schema: ntsReleaseShipmentSchema,
        showPickupFields: false,
        showDeliveryFields: true,
      };
    default:
      throw new Error('unrecognized shipment type');
  }
}

class MtoShipmentForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hasDeliveryAddress: get(props.mtoShipment, 'destinationAddress', false),
      useCurrentResidence: false,
      initialValues: {
        useCurrentResidence: false,
        customerRemarks: '',
        pickup: {
          address: {
            street_address_1: '',
            street_address_2: '',
            city: '',
            state: '',
            postal_code: '',
          },
          agent: {
            firstName: '',
            lastName: '',
            email: '',
            phone: '',
          },
        },
        delivery: {
          address: {
            street_address_1: '',
            street_address_2: '',
            city: '',
            state: '',
            postal_code: '',
          },
          agent: {
            firstName: '',
            lastName: '',
            email: '',
            phone: '',
          },
        },
      },
    };
  }

  componentDidMount() {
    const { showLoggedInUser } = this.props;
    showLoggedInUser();
  }

  handleChangeHasDeliveryAddress = () => {
    this.setState((prevState) => {
      return { hasDeliveryAddress: !prevState.hasDeliveryAddress };
    });
  };

  // Use current residence
  handleUseCurrentResidenceChange = (currentValues) => {
    const { initialValues } = this.state;
    const { currentResidence, match, mtoShipment } = this.props;
    this.setState(
      (state) => ({ useCurrentResidence: !state.useCurrentResidence }),
      () => {
        const { pickup } = currentValues;
        const { useCurrentResidence } = this.state;
        if (useCurrentResidence) {
          pickup.address = {
            street_address_1: currentResidence.street_address_1,
            street_address_2: currentResidence.street_address_2,
            city: currentResidence.city,
            state: currentResidence.state,
            postal_code: currentResidence.postal_code,
          };
        } else if (match.params.moveId === initialValues.moveTaskOrderID) {
          pickup.address = {
            street_address_1: mtoShipment.pickupAddress.street_address_1,
            street_address_2: mtoShipment.pickupAddress.street_address_2,
            city: mtoShipment.pickupAddress.city,
            state: mtoShipment.pickupAddress.state,
            postal_code: mtoShipment.pickupAddress.postal_code,
          };
        } else {
          pickup.address = {
            street_address_1: '',
            street_address_2: '',
            city: '',
            state: '',
            postal_code: '',
          };
        }

        // eslint-disable-next-line react/destructuring-assignment
        this.setState({
          initialValues: {
            ...initialValues,
            ...currentValues,
            pickup,
          },
        });
      },
    );
  };

  submitMTOShipment = ({ pickup, delivery, customerRemarks }) => {
    const { createMTOShipment, match, selectedMoveType } = this.props;
    const { moveId } = match.params;

    const pendingMtoShipment = formatMtoShipment({
      shipmentType: selectedMoveType,
      moveId,
      customerRemarks,
      pickup,
      delivery,
    });

    createMTOShipment(pendingMtoShipment);
  };

  render() {
    // TODO: replace minimal styling with actual styling during UI phase
    const { pageKey, pageList, match, history, newDutyStationAddress, selectedMoveType } = this.props;
    const { useCurrentResidence, hasDeliveryAddress, initialValues } = this.state;
    const fieldsetClasses = 'margin-top-2';
    const options = getShipmentOptions(selectedMoveType);

    return (
      <Formik
        initialValues={initialValues}
        enableReinitialize
        validateOnBlur
        validateOnChange
        validationSchema={options.schema}
      >
        {({ values, dirty, isValid }) => (
          <WizardPage
            canMoveNext={dirty && isValid}
            match={match}
            pageKey={pageKey}
            pageList={pageList}
            push={history.push}
            handleSubmit={() => this.submitMTOShipment(values, dirty)}
          >
            <h1>Now lets arrange details for the professional movers</h1>
            <Form className={styles.HHGDetailsForm}>
              {options.showPickupFields && (
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
                          onChange={() => this.handleUseCurrentResidenceChange(values)}
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
              {options.showDeliveryFields && (
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
                        onChange={this.handleChangeHasDeliveryAddress}
                        checked={hasDeliveryAddress}
                      />
                      <Radio
                        id="no-delivery-address"
                        label="No"
                        name="hasDeliveryAddress"
                        onChange={this.handleChangeHasDeliveryAddress}
                        checked={!hasDeliveryAddress}
                      />
                    </div>
                    {hasDeliveryAddress ? (
                      <AddressFields name="delivery.address" values={values.delivery.address} />
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
          </WizardPage>
        )}
      </Formik>
    );
  }
}

MtoShipmentForm.propTypes = {
  match: MatchShape,
  history: HistoryShape,
  pageList: PageListShape,
  pageKey: PageKeyShape,
  createMTOShipment: func.isRequired,
  showLoggedInUser: func.isRequired,
  currentResidence: AddressShape.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  selectedMoveType: string.isRequired,
  mtoShipment: HhgShipmentShape,
};

MtoShipmentForm.defaultProps = {
  pageList: [],
  pageKey: '',
  match: { isExact: false, params: { moveID: '' } },
  history: { goBack: () => {}, push: () => {} },
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

export { MtoShipmentForm as MtoShipmentFormComponent };
export default connect(mapStateToProps, mapDispatchToProps)(MtoShipmentForm);
