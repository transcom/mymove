import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Label, Button, Textarea } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';
import { useNavigate } from 'react-router';

import styles from './CreateSITExtensionRequestForm.module.scss';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { DropdownInput } from 'components/form/fields/DropdownInput';
import { ShipmentShape } from 'types/shipment';
import { sitExtensionReasons } from 'constants/sitExtensions';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { dropdownInputOptions } from 'utils/formatters';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import RequiredAsterisk, { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const createSITExtensionRequestValidationSchema = Yup.object().shape({
  requestReason: Yup.string().required('Required'),
  requestedDays: Yup.number().required('Required'),
  contractorRemarks: Yup.string().required('Required'),
});

const CreateSITExtensionRequestForm = ({ shipment, submission }) => {
  const navigate = useNavigate();
  const initialValues = {
    modelType: 'CreateSITExtension',
    requestReason: '',
    requestedDays: '',
    contractorRemarks: '',
    mtoShipmentID: shipment.id,
  };

  const onSubmit = (values) => {
    const { requestReason, requestedDays, contractorRemarks, mtoShipmentID } = values;

    const body = {
      requestReason,
      requestedDays: Number(requestedDays),
      contractorRemarks,
    };
    submission({ mtoShipmentID, body });
  };

  return (
    <SectionWrapper>
      <Formik
        initialValues={initialValues}
        validationSchema={createSITExtensionRequestValidationSchema}
        onSubmit={onSubmit}
      >
        <Form data-testid="CreateSITExtensionRequestForm">
          <div className={styles.container}>
            <input type="hidden" name="mtoShipmentID" />
            {requiredAsteriskMessage}
            <DropdownInput
              label="Request Reason"
              name="requestReason"
              id="requestReason"
              showRequiredAsterisk
              required
              options={dropdownInputOptions(sitExtensionReasons)}
            />
            <MaskedTextField
              data-testid="requestedDays"
              name="requestedDays"
              label="Requested Days"
              id="requestedDays"
              signed={false}
              mask={Number}
              scale={0}
              thousandsSeparator=","
              lazy={false}
              showRequiredAsterisk
              required
            />
            <Label htmlFor="contractorRemarksInput" required>
              <span required>
                Contractor Remarks <RequiredAsterisk />
              </span>
            </Label>
            <Field
              id="contractorRemarksInput"
              name="contractorRemarks"
              as={Textarea}
              required
              className={`${formStyles.remarks}`}
            />
            <div className={styles.buttonGroup}>
              <Button secondary onClick={() => navigate(-1)}>
                Back
              </Button>
              <Button type="submit">Request SIT Extension</Button>
            </div>
          </div>
        </Form>
      </Formik>
    </SectionWrapper>
  );
};

CreateSITExtensionRequestForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default CreateSITExtensionRequestForm;
