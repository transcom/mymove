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
import { HhgShipmentShape, WizardPageShape } from 'types/customerShapes';
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
        displayName: 'HHG',
      };
    case SHIPMENT_OPTIONS.NTS:
      return {
        schema: ntsShipmentSchema,
        showPickupFields: true,
        showDeliveryFields: false,
        displayName: 'NTS',
      };
    case SHIPMENT_OPTIONS.NTSR:
      return {
        schema: ntsReleaseShipmentSchema,
        showPickupFields: false,
        showDeliveryFields: true,
        displayName: 'NTS-R',
      };
    default:
      throw new Error('unrecognized shipment type');
  }
}

class MtoShipmentForm extends Component {
  constructor(props) {
    super(props);
    const hasDeliveryAddress = get(props.mtoShipment, 'destinationAddress', false);
    this.state = {
      hasDeliveryAddress,
      useCurrentResidence: false,
      initialValues: {
        pickup: {
          address: {},
          agent: {},
        },
        delivery: {
          address: {},
          agent: {},
        },
      },
    };
  }

  componentDidMount() {
    const { showLoggedInUser, isEditPage } = this.props;
    showLoggedInUser();

    // If refreshing edit page, need to handle mtoShipment populating from a promise
    if (isEditPage && mtoShipment.id) {
      this.setInitialState(mtoShipment);
    }
  }
  
  componentDidUpdate(prevProps) {
    const { mtoShipment, isEditPage } = this.props;

    // If refreshing edit page, need to handle mtoShipment populating from a promise
    if (isEditPage && mtoShipment.id && prevProps.mtoShipment.id !== mtoShipment.id) {
      this.setInitialEditState(mtoShipment);
    }
  }

  // Use current residence
  handleUseCurrentResidenceChange = (currentValues) => {
    const { initialValues } = this.state;
    const { currentResidence, wizardPage, mtoShipment } = this.props;
    this.setState(
      (state) => ({ useCurrentResidence: !state.useCurrentResidence }),
      () => {
        // eslint-disable-next-line react/destructuring-assignment
        if (this.state.useCurrentResidence) {
          this.setState({
            initialValues: {
              ...initialValues,
              ...currentValues,
              pickup: {
                address: {
                  street_address_1: currentResidence.street_address_1,
                  street_address_2: currentResidence.street_address_2,
                  city: currentResidence.city,
                  state: currentResidence.state,
                  postal_code: currentResidence.postal_code,
                },
              },
            },
          });
        } else {
          // eslint-disable-next-line no-lonely-if
          if (wizardPage.match.params.moveId === initialValues.moveTaskOrderID) {
            this.setState({
              initialValues: {
                ...initialValues,
                ...currentValues,
                pickup: {
                  address: {
                    street_address_1: mtoShipment.pickupAddress.street_address_1,
                    street_address_2: mtoShipment.pickupAddress.street_address_2,
                    city: mtoShipment.pickupAddress.city,
                    state: mtoShipment.pickupAddress.state,
                    postal_code: mtoShipment.pickupAddress.postal_code,
                  },
                },
              },
            });
          } else {
            this.setState({
              initialValues: {
                ...initialValues,
                ...currentValues,
                pickup: {
                  address: {
                    street_address_1: '',
                    street_address_2: '',
                    city: '',
                    state: '',
                    postal_code: '',
                  },
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

  submitMTOShipment = ({ pickup, delivery, customerRemarks }) => {
    const { createMTOShipment, updateMTOShipment, wizardPage, selectedMoveType } = this.props;
    const { moveId } = wizardPage.match.params;
    
    const pendingMtoShipment = formatMtoShipment({
      shipmentType: selectedMoveType,
      moveId,
      customerRemarks,
      pickup,
      delivery,
    });

    if (isEditPage) {
      updateMTOShipment(mtoShipment.id, pendingMtoShipment, mtoShipment.eTag).then(() => {
        wizardPage.history.goBack();
      });
    } else {
      createMTOShipment(pendingMtoShipment);
    }
  };
  
  // TODO: finish updating to match new initialState structure
  setInitialEditState = (mtoShipment) => {
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
    // for existing mtoShipment, reshape agents from array of objects to key/object for proper handling
    const { agents } = mtoShipment;
    const formattedMTOShipment = { ...mtoShipment };
    if (agents) {
      const receivingAgent = agents.find((agent) => agent.agentType === 'RECEIVING_AGENT');
      const releasingAgent = agents.find((agent) => agent.agentType === 'RELEASING_AGENT');

      // Remove dashes from agent phones for expected form phone format
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
    const { wizardPage, newDutyStationAddress, displayOptions } = this.props;
    const { pageKey, pageList, match, history } = wizardPage;
    const { hasDeliveryAddress, useCurrentResidence, initialValues } = this.state;
    const fieldsetClasses = 'margin-top-2';
    
    const editForm = (
      <div className="grid-container">
        <Formik
          initialValues={initialValues}
          enableReinitialize
          validateOnBlur
          validateOnChange
          validationSchema={displayOptions.schema}
        >
          {({ values, dirty, isValid,  isSubmitting, handleChange }) => (
            <MtoShipmentInnerForm
              {...this.props}
            />
          )}
        </Formik>
      </div>);
    
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
              <MtoShipmentInnerForm
              {...this.props}
            />
            </WizardPage>
          )}
        </Formik>
      </div>);

    return isEditPage ? editForm : createForm;
  }
}

MtoShipmentForm.propTypes = {
  wizardPage: WizardPageShape,
  createMTOShipment: func.isRequired,
  showLoggedInUser: func.isRequired,
  isEditPage: bool.isRequired,
  currentResidence: AddressShape.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  selectedMoveType: string.isRequired,
  mtoShipment: HhgShipmentShape,
};

MtoShipmentForm.defaultProps = {
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
    mtoShipment: selectMTOShipmentForMTO(state, ownProps.wizardPage.match.params.moveId),
    currentResidence: get(selectServiceMemberFromLoggedInUser(state), 'residential_address', {}),
    newDutyStationAddress: get(orders, 'new_duty_station.address', {}),
    displayOptions: getShipmentOptions(ownProps.selectedMoveType),
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
