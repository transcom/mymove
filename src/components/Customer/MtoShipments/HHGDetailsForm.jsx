import React, { Component } from 'react';
import { func } from 'prop-types';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { Fieldset } from '@trussworks/react-uswds';

import styles from './HHGDetailsForm.module.scss';
import { RequiredPlaceSchema, OptionalPlaceSchema } from './validationSchemas';
import { HhgShipmentShape, wizardPageShape } from './propShapes';
import { formatMtoShipment } from './utils';
import { PickupFields } from './PickupFields';
import { DeliveryFields } from './DeliveryFields';

import { TextInput } from 'components/form/fields';
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
import { AddressShape, SimpleAddressShape } from 'types/address';

const HHGDetailsFormSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  delivery: OptionalPlaceSchema,
  customerRemarks: Yup.string(),
});

class HHGDetailsForm extends Component {
  constructor(props) {
    super(props);
    const hasDeliveryAddress = get(props.mtoShipment, 'destinationAddress', false);
    this.state = {
      hasDeliveryAddress,
      useCurrentResidence: false,
      initialValues: {},
    };
  }

  componentDidMount() {
    const { showLoggedInUser } = this.props;
    showLoggedInUser();
  }

  submitMTOShipment = ({ pickup, delivery, customerRemarks }) => {
    const { createMTOShipment, match } = this.props;
    const { hasDeliveryAddress } = this.state;
    const { moveId } = match.params;

    const pendingMtoShipment = formatMtoShipment({
      moveId,
      pickup,
      customerRemarks,
      shipmentType: SHIPMENT_OPTIONS.HHG,
      delivery: hasDeliveryAddress ? delivery : undefined,
    });

    createMTOShipment(pendingMtoShipment);
  };

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
        // eslint-disable-next-line react/destructuring-assignment
        if (this.state.useCurrentResidence) {
          this.setState({
            initialValues: {
              ...initialValues,
              ...currentValues,
              pickupAddress: {
                street_address_1: currentResidence.street_address_1,
                street_address_2: currentResidence.street_address_2,
                city: currentResidence.city,
                state: currentResidence.state,
                postal_code: currentResidence.postal_code,
              },
            },
          });
        } else {
          // eslint-disable-next-line no-lonely-if
          if (match.params.moveId === initialValues.moveTaskOrderID) {
            this.setState({
              initialValues: {
                ...initialValues,
                ...currentValues,
                pickupAddress: {
                  street_address_1: mtoShipment.pickupAddress.street_address_1,
                  street_address_2: mtoShipment.pickupAddress.street_address_2,
                  city: mtoShipment.pickupAddress.city,
                  state: mtoShipment.pickupAddress.state,
                  postal_code: mtoShipment.pickupAddress.postal_code,
                },
              },
            });
          } else {
            this.setState({
              initialValues: {
                ...initialValues,
                ...currentValues,
                pickupAddress: {
                  street_address_1: '',
                  street_address_2: '',
                  city: '',
                  state: '',
                  postal_code: '',
                },
              },
            });
          }
        }
      },
    );
  };

  render() {
    // TODO: replace minimal styling with actual styling during UI phase
    const { pageKey, pageList, match, push, newDutyStationAddress } = this.props;
    const { hasDeliveryAddress, useCurrentResidence, initialValues } = this.state;
    const fieldsetClasses = 'margin-top-2';
    return (
      <Formik
        initialValues={initialValues}
        enableReinitialize
        validateOnBlur
        validateOnChange
        validationSchema={HHGDetailsFormSchema}
      >
        {({ values, dirty, isValid }) => (
          <WizardPage
            canMoveNext={dirty && isValid}
            match={match}
            pageKey={pageKey}
            pageList={pageList}
            push={push}
            handleSubmit={() => this.submitMTOShipment(values, dirty)}
          >
            <h1>Now lets arrange details for the professional movers</h1>
            <Form className={styles.HHGDetailsForm}>
              <PickupFields
                fieldsetClasses={fieldsetClasses}
                useCurrentResidence={useCurrentResidence}
                onCurrentResidenceChange={this.handleUseCurrentResidenceChange}
                values={values.pickup}
              />
              <DeliveryFields
                fieldsetClasses={fieldsetClasses}
                newDutyStationAddress={newDutyStationAddress}
                hasDeliveryAddress={hasDeliveryAddress}
                onHasAddressChange={this.handleChangeHasDeliveryAddress}
                values={values.delivery}
              />
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

HHGDetailsForm.propTypes = {
  wizardPageShape,
  createMTOShipment: func.isRequired,
  showLoggedInUser: func.isRequired,
  currentResidence: AddressShape.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  mtoShipment: HhgShipmentShape,
};

HHGDetailsForm.defaultProps = {
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

export { HHGDetailsForm as HHGDetailsFormComponent };
export default connect(mapStateToProps, mapDispatchToProps)(HHGDetailsForm);
