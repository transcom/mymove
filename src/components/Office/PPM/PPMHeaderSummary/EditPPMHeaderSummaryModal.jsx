import { React, useState } from 'react';
import { Formik } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button } from '@trussworks/react-uswds';

import styles from './EditPPMHeaderSummaryModal.module.scss';

import { formatCentsTruncateWhole } from 'utils/formatters';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import { DatePickerInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { requiredAddressSchema } from 'utils/validation';

const EditPPMHeaderSummaryModal = ({ sectionType, sectionInfo, onClose, onSubmit, editItemName }) => {
  const { actualMoveDate, advanceAmountReceived, pickupAddressObj, destinationAddressObj } = sectionInfo;
  let title = 'Edit';
  if (sectionType === 'shipmentInfo') {
    title = 'Edit Shipment Info';
  } else if (sectionType === 'incentives') {
    title = 'Edit Incentives/Costs';
  }
  const initialValues = {
    editItemName,
    actualMoveDate: actualMoveDate || '',
    advanceAmountReceived: formatCentsTruncateWhole(advanceAmountReceived).replace(/,/g, ''),
    pickupAddress: pickupAddressObj,
    destinationAddress: destinationAddressObj,
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
  const [isLookupErrorVisible, setIsLookupErrorVisible] = useState(false);

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.EditPPMHeaderSummaryModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle className={styles.ModalTitle}>
            <h3>{title}</h3>
          </ModalTitle>
          <Formik validationSchema={validationSchema} initialValues={initialValues} onSubmit={onSubmit}>
            {({ isValid, handleChange, setFieldTouched, values, setValues }) => {
              const handlePickupZipCityChange = (value) => {
                setValues(
                  {
                    ...values,
                    pickupAddress: {
                      ...values.pickupAddress,
                      city: value.city,
                      state: value.state ? value.state : '',
                      county: value.county,
                      postalCode: value.postalCode,
                    },
                  },
                  { shouldValidate: true },
                );

                if (!value.city || !value.state || !value.county || !value.postalCode) {
                  setIsLookupErrorVisible(true);
                } else {
                  setIsLookupErrorVisible(false);
                }
              };
              const handleDestinationZipCityChange = (value) => {
                setValues(
                  {
                    ...values,
                    destinationAddress: {
                      ...values.destinationAddress,
                      city: value.city,
                      state: value.state ? value.state : '',
                      county: value.county,
                      postalCode: value.postalCode,
                    },
                  },
                  { shouldValidate: true },
                );

                if (!value.city || !value.state || !value.county || !value.postalCode) {
                  setIsLookupErrorVisible(true);
                } else {
                  setIsLookupErrorVisible(false);
                }
              };
              return (
                <Form>
                  <div>
                    {editItemName === 'actualMoveDate' && (
                      <DatePickerInput
                        name="actualMoveDate"
                        label="Actual move start date"
                        id="actualMoveDate"
                        disabledDays={{ after: new Date() }}
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
                        legend="Pickup Address"
                        className={styles.AddressFieldSet}
                        formikFunctionsToValidatePostalCodeOnChange={{ handleChange, setFieldTouched }}
                        zipCityEnabled
                        zipCityError={isLookupErrorVisible}
                        handleLocationChange={handlePickupZipCityChange}
                      />
                    )}
                    {editItemName === 'destinationAddress' && (
                      <AddressFields
                        name="destinationAddress"
                        legend="Destination Address"
                        className={styles.AddressFieldSet}
                        formikFunctionsToValidatePostalCodeOnChange={{ handleChange, setFieldTouched }}
                        zipCityEnabled
                        zipCityError={isLookupErrorVisible}
                        handleLocationChange={handleDestinationZipCityChange}
                      />
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
