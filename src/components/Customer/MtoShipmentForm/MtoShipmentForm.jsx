/* eslint-disable react/jsx-props-no-spreading */
import React, { Component } from 'react';
import { bool, string, func, shape, number } from 'prop-types';
import { Formik } from 'formik';

import getShipmentOptions from './getShipmentOptions';
import MtoShipmentFormFields from './MtoShipmentFormFields';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { WizardPage } from 'shared/WizardPage/index';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { HhgShipmentShape, MatchShape, HistoryShape, PageKeyShape, PageListShape } from 'types/customerShapes';
import { formatMtoShipmentForAPI, formatMtoShipmentForDisplay } from 'utils/formatMtoShipment';
import { createMTOShipment, patchMTOShipment, getResponseError } from 'services/internalApi';
import Alert from 'shared/Alert';

class MtoShipmentForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMessage: null,
    };
  }

  submitMTOShipment = ({ shipmentType, pickup, hasDeliveryAddress, delivery, customerRemarks }) => {
    const { history, match, selectedMoveType, isCreatePage, mtoShipment, updateMTOShipment } = this.props;
    const { moveId } = match.params;

    const deliveryDetails = delivery;
    if (hasDeliveryAddress === 'no') {
      delete deliveryDetails.address;
    }

    const pendingMtoShipment = formatMtoShipmentForAPI({
      shipmentType: shipmentType || selectedMoveType,
      moveId,
      customerRemarks,
      pickup,
      delivery: deliveryDetails,
    });

    if (isCreatePage) {
      createMTOShipment(pendingMtoShipment)
        .then((response) => {
          updateMTOShipment(response);
        })
        .catch((e) => {
          const { response } = e;
          const errorMessage = getResponseError(response, 'failed to create MTO shipment due to server error');

          this.setState({ errorMessage });
        });
    } else {
      patchMTOShipment(mtoShipment.id, pendingMtoShipment, mtoShipment.eTag)
        .then((response) => {
          updateMTOShipment(response);
          history.push(`/moves/${moveId}/review`);
        })
        .catch((e) => {
          const { response } = e;
          const errorMessage = getResponseError(response, 'failed to update MTO shipment due to server error');

          this.setState({ errorMessage });
        });
    }
  };

  getShipmentNumber = () => {
    // TODO - this is not supported by IE11, shipment number should be calculable from Redux anyways
    // we should fix this also b/c it doesn't display correctly in storybook
    const { search } = window.location;
    const params = new URLSearchParams(search);
    const shipmentNumber = params.get('shipmentNumber');
    return shipmentNumber;
  };

  render() {
    const {
      pageKey,
      pageList,
      match,
      history,
      newDutyStationAddress,
      selectedMoveType,
      isCreatePage,
      mtoShipment,
      serviceMember,
      currentResidence,
    } = this.props;
    const { errorMessage } = this.state;

    const shipmentType = selectedMoveType || mtoShipment.shipmentType;
    const { showDeliveryFields, showPickupFields, schema } = getShipmentOptions(shipmentType);
    const initialValues = formatMtoShipmentForDisplay(isCreatePage ? {} : mtoShipment);

    const commonFormProps = {
      isCreatePage,
      pageKey,
      pageList,
      match,
      history,
      newDutyStationAddress,
      serviceMember,
      showPickupFields,
      showDeliveryFields,
      shipmentType,
      shipmentNumber: shipmentType === SHIPMENT_OPTIONS.HHG ? this.getShipmentNumber() : null,
    };

    return (
      <Formik
        initialValues={initialValues}
        enableReinitialize
        validateOnBlur
        validateOnChange
        validationSchema={schema}
      >
        {({ values, dirty, isValid, isSubmitting, setValues }) => {
          const handleUseCurrentResidenceChange = (e) => {
            const { checked } = e.target;
            if (checked) {
              // use current residence
              setValues({
                ...values,
                pickup: {
                  ...values.pickup,
                  address: currentResidence,
                },
              });
            } else if (match.params.moveId === mtoShipment?.moveTaskOrderId) {
              // TODO - what is the purpose of this check?
              // Revert address
              setValues({
                ...values,
                pickup: {
                  ...values.pickup,
                  address: mtoShipment.pickupAddress,
                },
              });
            } else {
              // Revert address
              setValues({
                ...values,
                pickup: {
                  ...values.pickup,
                  address: {
                    street_address_1: '',
                    street_address_2: '',
                    city: '',
                    state: '',
                    postal_code: '',
                  },
                },
              });
            }
          };

          if (isCreatePage) {
            // return MTO Shipment form in the wizard
            return (
              <WizardPage
                canMoveNext={dirty && isValid}
                match={match}
                pageKey={pageKey}
                pageList={pageList}
                push={history.push}
                handleSubmit={() => this.submitMTOShipment(values, dirty)}
                error={errorMessage}
              >
                <MtoShipmentFormFields
                  {...commonFormProps}
                  values={values}
                  onUseCurrentResidenceChange={handleUseCurrentResidenceChange}
                  submitHandler={this.submitMTOShipment}
                />
              </WizardPage>
            );
          }

          return (
            <div className="grid-container usa-prose">
              <div className="grid-row">
                <div className="grid-col">
                  {errorMessage && (
                    <div className="usa-width-one-whole error-message">
                      <Alert type="error" heading="An error occurred">
                        {errorMessage}
                      </Alert>
                    </div>
                  )}
                  <MtoShipmentFormFields
                    {...commonFormProps}
                    values={values}
                    onUseCurrentResidenceChange={handleUseCurrentResidenceChange}
                    submitHandler={this.submitMTOShipment}
                    dirty={dirty}
                    isValid={isValid}
                    isSubmitting={isSubmitting}
                  />
                </div>
              </div>
            </div>
          );
        }}
      </Formik>
    );
  }
}

MtoShipmentForm.propTypes = {
  match: MatchShape,
  history: HistoryShape,
  pageList: PageListShape,
  pageKey: PageKeyShape,
  updateMTOShipment: func.isRequired,
  isCreatePage: bool,
  currentResidence: AddressShape.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  selectedMoveType: string.isRequired,
  mtoShipment: HhgShipmentShape,
  serviceMember: shape({
    weight_allotment: shape({
      total_weight_self: number,
    }),
  }).isRequired,
};

MtoShipmentForm.defaultProps = {
  isCreatePage: false,
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

export default MtoShipmentForm;
