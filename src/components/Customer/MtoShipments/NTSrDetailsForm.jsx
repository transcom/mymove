import React, { Component } from 'react';
import { func } from 'prop-types';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { Fieldset } from '@trussworks/react-uswds';

import styles from './MtoShipmentFormStyles.module.scss';
import { RequiredPlaceSchema } from './validationSchemas';
import { NtsrShipmentShape, WizardPageShape } from './propShapes';
import { DeliveryFields } from './FormGroups/DeliveryFields';

import { TextInput } from 'components/form/fields';
import { Form } from 'components/form/Form';
import {
  selectMTOShipmentForMTO,
  createMTOShipment as createMTOShipmentAction,
  updateMTOShipment as updateMTOShipmentAction,
} from 'shared/Entities/modules/mtoShipments';
import { selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { WizardPage } from 'shared/WizardPage';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { SimpleAddressShape } from 'types/address';
import { formatMtoShipment } from 'utils/formatMtoShipment';

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
    const { createMTOShipment, wizardPage } = this.props;
    const { hasDeliveryAddress } = this.state;
    const { moveId } = wizardPage.match.params;

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
    const { wizardPage, newDutyStationAddress } = this.props;
    const { pageKey, pageList, match, history } = wizardPage;
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
            push={history.push}
            handleSubmit={() => this.submitMTOShipment(values, dirty)}
          >
            <h1>Now lets arrange details for the professional movers</h1>
            <Form className={styles.HHGDetailsForm}>
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

NTSrDetailsForm.propTypes = {
  wizardPage: WizardPageShape,
  createMTOShipment: func.isRequired,
  showLoggedInUser: func.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  mtoShipment: NtsrShipmentShape,
};

NTSrDetailsForm.defaultProps = {
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
