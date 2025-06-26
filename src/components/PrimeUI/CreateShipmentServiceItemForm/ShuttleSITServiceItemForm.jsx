import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import formStyles from 'styles/form.module.scss';
import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { Form } from 'components/form/Form';
import { ShipmentShape } from 'types/shipment';
import { DropdownInput } from 'components/form/fields';
import { domesticShuttleServiceItemCodeOptions, createServiceItemModelTypes } from 'constants/prime';

const shuttleSITValidationSchema = Yup.object().shape({
  reServiceCode: Yup.string().required('Required'),
  reason: Yup.string().required('Required'),
});

const ShuttleSITServiceItemForm = ({ shipment, submission, handleCancel }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: createServiceItemModelTypes.MTOServiceItemDomesticShuttle,
    reason: '',
    estimatedWeight: null,
    actualWeight: null,
  };

  const onSubmit = (values) => {
    const { estimatedWeight, actualWeight, ...otherFields } = values;
    const body = {
      estimatedWeight: Number.parseInt(estimatedWeight, 10),
      actualWeight: Number.parseInt(actualWeight, 10),
      ...otherFields,
    };
    submission({ body });
  };

  return (
    <Formik initialValues={initialValues} validationSchema={shuttleSITValidationSchema} onSubmit={onSubmit}>
      <Form data-testid="shuttleSITServiceItemForm" className={formStyles.form}>
        <DropdownInput
          label="Service item code"
          name="reServiceCode"
          id="reServiceCode"
          required
          options={domesticShuttleServiceItemCodeOptions}
          showRequiredAsterisk
        />
        <TextField name="reason" id="reason" label="Reason" showRequiredAsterisk required />
        <MaskedTextField
          data-testid="estimatedWeightInput"
          name="estimatedWeight"
          label="Estimated weight (lbs)"
          id="estimatedWeightInput"
          mask={Number}
          scale={0}
          thousandsSeparator=","
          lazy={false}
        />
        <MaskedTextField
          data-testid="actualWeightInput"
          name="actualWeight"
          label="Actual weight (lbs)"
          id="actualWeightInput"
          mask={Number}
          scale={0}
          thousandsSeparator=","
          lazy={false}
        />
        <div className={formStyles.formActions}>
          <Button type="button" secondary onClick={handleCancel}>
            Cancel
          </Button>
          <Button type="submit">Create service item</Button>
        </div>
      </Form>
    </Formik>
  );
};

ShuttleSITServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default ShuttleSITServiceItemForm;
