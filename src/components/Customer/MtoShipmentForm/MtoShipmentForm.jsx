/* eslint-disable react/jsx-props-no-spreading */
import React, { Component } from 'react';
import { bool, string, func, shape } from 'prop-types';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { Formik } from 'formik';

import { getShipmentOptions } from './getShipmentOptions';
import MtoShipmentFormFields from './MtoShipmentFormFields';

import {
  selectMTOShipmentById,
  createMTOShipment as createMTOShipmentAction,
  updateMTOShipment as updateMTOShipmentAction,
} from 'shared/Entities/modules/mtoShipments';
import { selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { WizardPage } from 'shared/WizardPage';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { HhgShipmentShape, MatchShape, HistoryShape, PageKeyShape, PageListShape } from 'types/customerShapes';
import { formatMtoShipmentForAPI, formatMtoShipmentForDisplay } from 'utils/formatMtoShipment';

class MtoShipmentForm extends Component {
  componentDidMount() {
    const { showLoggedInUser } = this.props;
    showLoggedInUser();
    // TODO - move this to the parent component instead

    // TODO: confirm this block should exist
    // If refreshing edit page, need to handle mtoShipment populating from a promise
    /*
    if (!isCreatePage && mtoShipment.id) {
      this.setState({
        initialValues,
        hasDeliveryAddress: initialValues.hasDeliveryAddress,
      });
    } */
  }

  submitMTOShipment = ({ shipmentType, pickup, delivery, customerRemarks }) => {
    const {
      createMTOShipment,
      updateMTOShipment,
      history,
      match,
      selectedMoveType,
      isCreatePage,
      mtoShipment,
    } = this.props;
    const { moveId } = match.params;

    const pendingMtoShipment = formatMtoShipmentForAPI({
      shipmentType: shipmentType || selectedMoveType,
      moveId,
      customerRemarks,
      pickup,
      delivery,
    });

    if (isCreatePage) {
      createMTOShipment(pendingMtoShipment);
    } else {
      updateMTOShipment(mtoShipment.id, pendingMtoShipment, mtoShipment.eTag).then(() => {
        history.goBack();
      });
    }
  };

  getShipmentNumber = () => {
    // TODO - fix
    const { search } = window.location;
    const params = new URLSearchParams(search);
    const shipmentNumber = params.get('shipmentNumber');
    return shipmentNumber;
  };

  render() {
    // TODO: replace minimal styling with actual styling during UI phase
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

    const displayOptions = getShipmentOptions(selectedMoveType || mtoShipment.shipmentType);
    const initialValues = formatMtoShipmentForDisplay(isCreatePage ? {} : mtoShipment);

    const commonFormProps = {
      isCreatePage,
      pageKey,
      pageList,
      match,
      history,
      newDutyStationAddress,
      displayOptions,
      serviceMember,
    };

    return (
      <div className="grid-container">
        <Formik
          initialValues={initialValues}
          enableReinitialize
          validateOnBlur
          validateOnChange
          validationSchema={displayOptions.schema}
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
              <MtoShipmentFormFields
                {...commonFormProps}
                values={values}
                onUseCurrentResidenceChange={handleUseCurrentResidenceChange}
                submitHandler={this.submitMTOShipment}
                dirty={dirty}
                isValid={isValid}
                isSubmitting={isSubmitting}
              />
            );
          }}
        </Formik>
      </div>
    );
  }
}

MtoShipmentForm.propTypes = {
  match: MatchShape,
  history: HistoryShape,
  pageList: PageListShape,
  pageKey: PageKeyShape,
  createMTOShipment: func.isRequired,
  updateMTOShipment: func.isRequired,
  showLoggedInUser: func.isRequired,
  isCreatePage: bool,
  currentResidence: AddressShape.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  selectedMoveType: string.isRequired,
  mtoShipment: HhgShipmentShape,
  serviceMember: shape({
    weight_allotment: shape({
      total_weight_self: string,
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

const mapStateToProps = (state, ownProps) => {
  const orders = selectActiveOrLatestOrdersFromEntities(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  const props = {
    serviceMember,
    mtoShipment: selectMTOShipmentById(state, ownProps.match.params.mtoShipmentId),
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
