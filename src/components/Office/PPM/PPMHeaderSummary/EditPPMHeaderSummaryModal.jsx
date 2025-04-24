import React, { useEffect, useState } from 'react';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, FormGroup, Radio, Label } from '@trussworks/react-uswds';

import styles from './EditPPMHeaderSummaryModal.module.scss';

import { formatCentsTruncateWhole, formatWeight } from 'utils/formatters';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import { DatePickerInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { requiredAddressSchema } from 'utils/validation';
import { FEATURE_FLAG_KEYS, getPPMTypeLabel, PPM_TYPES } from 'shared/constants';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const EditPPMHeaderSummaryModal = ({ sectionType, sectionInfo, onClose, onSubmit, editItemName, grade }) => {
  const [ppmSprFF, setPpmSprFF] = useState(false);
  const { ppmType, actualMoveDate, advanceAmountReceived, allowableWeight, pickupAddressObj, destinationAddressObj } =
    sectionInfo;
  let title = 'Edit';
  if (sectionType === 'shipmentInfo') {
    title = 'Edit Shipment Info';
  } else if (sectionType === 'incentives') {
    title = 'Edit Incentives/Costs';
  }
  const initialValues = {
    editItemName,
    ppmType: sectionInfo.ppmType,
    actualMoveDate: actualMoveDate || '',
    advanceAmountReceived: formatCentsTruncateWhole(advanceAmountReceived).replace(/,/g, ''),
    allowableWeight: formatWeight(allowableWeight),
    pickupAddress: pickupAddressObj,
    destinationAddress: destinationAddressObj,
    isActualExpenseReimbursement: sectionInfo.isActualExpenseReimbursement ? 'true' : 'false',
  };

  const validationSchema = Yup.object().shape({
    actualMoveDate: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .when('editItemName', {
        is: 'actualMoveDate',
        then: (schema) => schema.required('Required').max(new Date(), 'Date cannot be in the future'),
      }),
    advanceAmountReceived: Yup.number().when('editItemName', {
      is: 'advanceAmountReceived',
      then: (schema) => schema.required('Required'),
    }),
    pickupAddress: Yup.object().when('editItemName', {
      is: 'pickupAddress',
      then: () => requiredAddressSchema,
      otherwise: (schema) => schema,
    }),
    destinationAddress: Yup.object().when('editItemName', {
      is: 'destinationAddress',
      then: () => requiredAddressSchema,
      otherwise: (schema) => schema,
    }),
  });

  const weightValidationSchema = Yup.object().shape({
    allowableWeight: Yup.number().when('editItemName', {
      is: 'allowableWeight',
      then: (schema) => schema.required('Required').min(0, 'Allowable weight must be greater than or equal to zero'),
      otherwise: (schema) => schema,
    }),
  });

  useEffect(() => {
    const fetchData = async () => {
      setPpmSprFF(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.PPM_SPR));
    };
    fetchData();
  }, []);

  const isCivilian = grade === 'CIVILIAN_EMPLOYEE';

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.EditPPMHeaderSummaryModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle className={styles.ModalTitle}>
            <h3>{title}</h3>
          </ModalTitle>
          <Formik
            validationSchema={editItemName === 'allowableWeight' ? weightValidationSchema : validationSchema}
            initialValues={initialValues}
            onSubmit={onSubmit}
          >
            {({ isValid, handleChange, values, ...formikProps }) => {
              return (
                <Form>
                  <div>
                    {editItemName === 'actualMoveDate' && (
                      <DatePickerInput
                        name="actualMoveDate"
                        label="Actual move start date"
                        id="actualMoveDate"
                        disabledDays={{ after: new Date() }}
                        formikFunctionsToValidatePostalCodeOnChange={{
                          handleChange,
                          setFieldTouched: formikProps.setFieldTouched,
                        }}
                      />
                    )}
                    {editItemName === 'advanceAmountReceived' && (
                      <MaskedTextField
                        label="Advance received"
                        name="advanceAmountReceived"
                        id="advanceAmountReceived"
                        defaultValue="0"
                        mask={Number}
                        scale={0} // digits after point, 0 for integers
                        signed={false} // disallow negative
                        thousandsSeparator=","
                        lazy={false} // immediate masking evaluation
                        prefix="$"
                      />
                    )}
                    {editItemName === 'pickupAddress' && (
                      <AddressFields
                        name="pickupAddress"
                        legend={ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Shipped from Address' : 'Pickup Address'}
                        className={styles.AddressFieldSet}
                        formikProps={formikProps}
                      />
                    )}
                    {editItemName === 'destinationAddress' && (
                      <AddressFields
                        name="destinationAddress"
                        legend={ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Destination Address' : 'Delivery Address'}
                        className={styles.AddressFieldSet}
                        formikProps={formikProps}
                      />
                    )}
                    {editItemName === 'allowableWeight' && (
                      <MaskedTextField
                        label="Allowable Weight"
                        name="allowableWeight"
                        id="allowableWeight"
                        defaultValue="0"
                        mask={Number}
                        scale={0} // digits after point, 0 for integers
                        signed={false} // disallow negative
                        thousandsSeparator=","
                        lazy={false} // immediate masking evaluation
                        suffix="lbs"
                        data-testid="editAllowableWeightInput"
                      />
                    )}
                    {editItemName === 'expenseType' && (
                      <FormGroup>
                        <Label className={styles.Label} htmlFor="ppmType">
                          What is the PPM type?
                        </Label>
                        <Field
                          as={Radio}
                          id="isIncentiveBased"
                          label={getPPMTypeLabel(PPM_TYPES.INCENTIVE_BASED)}
                          name="ppmType"
                          value={PPM_TYPES.INCENTIVE_BASED}
                          disabled={isCivilian}
                          checked={values.ppmType === PPM_TYPES.INCENTIVE_BASED}
                          className={styles.buttonGroup}
                          data-testid="isIncentiveBased"
                        />
                        <Field
                          as={Radio}
                          id="isActualExpense"
                          label={getPPMTypeLabel(PPM_TYPES.ACTUAL_EXPENSE)}
                          name="ppmType"
                          value={PPM_TYPES.ACTUAL_EXPENSE}
                          checked={values.ppmType === PPM_TYPES.ACTUAL_EXPENSE}
                          className={styles.buttonGroup}
                          data-testid="isActualExpense"
                        />
                        {ppmSprFF && (
                          <Field
                            as={Radio}
                            id="isSmallPackage"
                            label={getPPMTypeLabel(PPM_TYPES.SMALL_PACKAGE)}
                            name="ppmType"
                            value={PPM_TYPES.SMALL_PACKAGE}
                            checked={values.ppmType === PPM_TYPES.SMALL_PACKAGE}
                            className={styles.buttonGroup}
                            data-testid="isSmallPackage"
                          />
                        )}
                      </FormGroup>
                    )}
                  </div>
                  <ModalActions>
                    <Button type="submit" disabled={!isValid}>
                      Save
                    </Button>
                    <Button
                      type="button"
                      onClick={() => onClose()}
                      data-testid="modalCancelButton"
                      outline
                      className={styles.CancelButton}
                    >
                      Cancel
                    </Button>
                  </ModalActions>
                </Form>
              );
            }}
          </Formik>
        </Modal>
      </ModalContainer>
    </div>
  );
};

EditPPMHeaderSummaryModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  editItemName: PropTypes.string.isRequired,
};

export default EditPPMHeaderSummaryModal;
