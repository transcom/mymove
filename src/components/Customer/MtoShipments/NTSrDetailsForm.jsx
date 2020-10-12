import React, { Component } from 'react';
import { arrayOf, string, bool, shape, func } from 'prop-types';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { Fieldset } from '@trussworks/react-uswds';

import { TextInput } from '../../form/fields';
import { Form } from '../../form/Form';

import styles from './HHGDetailsForm.module.scss';
import { RequiredPlaceSchema } from './formTypes';
import { simpleAddressShape, fullAddressShape, agentShape } from './propShapes';
import { formatMtoShipment } from './utils';
import { DeliveryDetails } from './DeliveryDetails';

import {
  selectMTOShipmentForMTO,
  createMTOShipment as createMTOShipmentAction,
  updateMTOShipment as updateMTOShipmentAction,
} from 'shared/Entities/modules/mtoShipments';
import { selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { WizardPage } from 'shared/WizardPage';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const NTSrDetailsFormSchema = Yup.object().shape({
  delivery: RequiredPlaceSchema,
  customerRemarks: Yup.string(),
});

class NTSrDetailsForm extends Component {
  constructor(props) {
    super(props);
    const hasDeliveryAddress = get(props.mtoShipment, 'destinationAddress', false);
    this.state = {
      hasDeliveryAddress,
      initialValues: {},
    };
  }

  componentDidMount() {
    const { showLoggedInUser } = this.props;
    showLoggedInUser();
  }

  submitMTOShipment = ({ delivery, customerRemarks }) => {
    const { createMTOShipment, match } = this.props;
    const { hasDeliveryAddress } = this.state;
    const { moveId } = match.params;

    const pendingMtoShipment = formatMtoShipment({
      moveId,
      customerRemarks,
      shipmentType: SHIPMENT_OPTIONS.NTSR,
      delivery: hasDeliveryAddress ? delivery : undefined,
    });

    createMTOShipment(pendingMtoShipment);
  };

  handleChangeHasDeliveryAddress = () => {
    this.setState((prevState) => {
      return { hasDeliveryAddress: !prevState.hasDeliveryAddress };
    });
  };

  render() {
    const { pageKey, pageList, match, push, newDutyStationAddress } = this.props;
    const { hasDeliveryAddress, initialValues } = this.state;
    const fieldsetClasses = 'margin-top-2';
    return (
      <Formik
        initialValues={initialValues}
        enableReinitialize
        validateOnBlur
        validateOnChange
        validationSchema={NTSrDetailsFormSchema}
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
            <Form className={styles.NTSrDetailsForm}>
              <DeliveryDetails
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

NTSrDetailsForm.propTypes = {
  pageKey: string.isRequired,
  pageList: arrayOf(string).isRequired,
  match: shape({
    isExact: bool.isRequired,
    params: shape({
      moveId: string.isRequired,
    }),
    path: string.isRequired,
    url: string.isRequired,
  }).isRequired,
  newDutyStationAddress: simpleAddressShape,
  createMTOShipment: func.isRequired,
  showLoggedInUser: func.isRequired,
  push: func.isRequired,
  mtoShipment: shape({
    agents: arrayOf(agentShape),
    customerRemarks: string,
    requestedDeliveryDate: string,
    destinationAddress: fullAddressShape,
  }),
};

NTSrDetailsForm.defaultProps = {
  newDutyStationAddress: {
    city: '',
    state: '',
    postal_code: '',
  },
  mtoShipment: {
    id: '',
    customerRemarks: '',
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
    newDutyStationAddress: get(orders, 'new_duty_station.address', {}),
  };
  return props;
};

const mapDispatchToProps = {
  createMTOShipment: createMTOShipmentAction,
  updateMTOShipment: updateMTOShipmentAction,
  showLoggedInUser: showLoggedInUserAction,
};

export { NTSrDetailsForm as NTSrDetailsFormComponent };
export default connect(mapStateToProps, mapDispatchToProps)(NTSrDetailsForm);
