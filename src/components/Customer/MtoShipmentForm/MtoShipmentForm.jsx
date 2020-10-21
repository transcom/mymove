import React, { Component } from 'react';
import { bool, string, func } from 'prop-types';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { Formik } from 'formik';

import { getShipmentOptions } from './getShipmentOptions';
import { MtoShipmentFormFields } from './MtoShipmentFormFields';

import {
  selectMTOShipmentForMTO,
  createMTOShipment as createMTOShipmentAction,
  updateMTOShipment as updateMTOShipmentAction,
} from 'shared/Entities/modules/mtoShipments';
import { selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { WizardPage } from 'shared/WizardPage';
import { AddressShape, SimpleAddressShape } from 'types/address';
import {
  HhgShipmentShape,
  MtoDisplayOptionsShape,
  MatchShape,
  HistoryShape,
  PageKeyShape,
  PageListShape,
} from 'types/customerShapes';
import { formatMtoShipment } from 'utils/formatMtoShipment';

/**
 * cleanAgentPhone removes dashes from agent phones for expected form phone format
 * (handling diff between expected FE and BE phone format)
 * @param {*} agent
 */
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

class MtoShipmentForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hasDeliveryAddress: get(props.mtoShipment, 'destinationAddress', false),
      useCurrentResidence: false,
      initialValues: {
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
    const { showLoggedInUser, isCreatePage, mtoShipment } = this.props;
    showLoggedInUser();

    // TODO: confirm this block should exist
    // If refreshing edit page, need to handle mtoShipment populating from a promise
    if (!isCreatePage && mtoShipment.id) {
      this.setInitialState(mtoShipment);
    }
  }

  componentDidUpdate(prevProps) {
    const { mtoShipment, isCreatePage } = this.props;

    // If refreshing edit page, need to handle mtoShipment populating from a promise
    if (!isCreatePage && mtoShipment.id && prevProps.mtoShipment.id !== mtoShipment.id) {
      this.setInitialEditState(mtoShipment);
    }
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

    const pendingMtoShipment = formatMtoShipment({
      shipmentType: selectedMoveType,
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

  // TODO: finish updating to match new initialState structure
  setInitialEditState = (mtoShipment) => {
    // for existing mtoShipment, reshape agents from array of objects to key/object for proper handling
    const { agents } = mtoShipment;
    const formattedMTOShipment = { ...mtoShipment };
    if (agents) {
      const receivingAgent = agents.find((agent) => agent.agentType === 'RECEIVING_AGENT');
      const releasingAgent = agents.find((agent) => agent.agentType === 'RELEASING_AGENT');

      if (receivingAgent) {
        const formattedAgent = cleanAgentPhone(receivingAgent);
        if (Object.keys(formattedAgent).length) {
          formattedMTOShipment.delivery.agent = { ...formattedAgent };
        }
      }
      if (releasingAgent) {
        const formattedAgent = cleanAgentPhone(releasingAgent);
        if (Object.keys(formattedAgent).length) {
          formattedMTOShipment.pickup.agent = { ...formattedAgent };
        }
      }
    }
    const hasDeliveryAddress = get(mtoShipment, 'destinationAddress', false);
    this.setState({ initialValues: formattedMTOShipment, hasDeliveryAddress });
  };

  getShipmentNumber = () => {
    const { search } = window.location;
    const params = new URLSearchParams(search);
    const shipmentNumber = params.get('shipmentNumber');
    return shipmentNumber;
  };

  render() {
    // TODO: replace minimal styling with actual styling during UI phase
    const { pageKey, pageList, match, history, newDutyStationAddress, selectedMoveType, isCreatePage } = this.props;
    const { useCurrentResidence, hasDeliveryAddress, initialValues } = this.state;
    const displayOptions = getShipmentOptions(selectedMoveType);
    const commonFormProps = {
      pageKey,
      pageList,
      match,
      history,
      newDutyStationAddress,
      selectedMoveType,
      isCreatePage,
    };

    const editForm = (
      <div className="grid-container">
        <Formik
          initialValues={initialValues}
          enableReinitialize
          validateOnBlur
          validateOnChange
          validationSchema={displayOptions.schema}
        >
          {({ values, dirty, isValid, isSubmitting }) => (
            <MtoShipmentFormFields
              {...commonFormProps}
              values={values}
              onHasDeliveryAddressChange={this.handleChangeHasDeliveryAddress}
              onUseCurrentResidenceChange={this.handleUseCurrentResidenceChange}
              submitHandler={this.submitMTOShipment}
              dirty={dirty}
              isValid={isValid}
              isSubmitting={isSubmitting}
            />
          )}
        </Formik>
      </div>
    );

    const createForm = (
      <div className="grid-container">
        <Formik
          initialValues={initialValues}
          enableReinitialize
          validateOnBlur
          validateOnChange
          validationSchema={displayOptions.schema}
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
              <MtoShipmentFormFields
                {...commonFormProps}
                values={values}
                onHasDeliveryAddressChange={this.handleChangeHasDeliveryAddress}
                onUseCurrentResidenceChange={this.handleUseCurrentResidenceChange}
                submitHandler={this.submitMTOShipment}
              />
            </WizardPage>
          )}
        </Formik>
      </div>
    );

    return isCreatePage ? createForm : editForm;
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
  displayOptions: MtoDisplayOptionsShape.isRequired,
  mtoShipment: HhgShipmentShape,
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
