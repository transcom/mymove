import React from 'react';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Label, TextInput, Fieldset, FormGroup, Grid } from '@trussworks/react-uswds';

import styles from './EditFacilityInfoModal.module.scss';

import { StorageFacilityAddressShape, StorageFacilityShape, ShipmentOptionsOneOf } from 'types/shipment';
import { Form } from 'components/form';
import formStyles from 'styles/form.module.scss';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { requiredAddressSchema, phoneSchema, emailSchema } from 'utils/validation';

const EditFacilityInfoModal = ({ onClose, onSubmit, storageFacility, storageFacilityAddress, shipmentType }) => {
  const editFacilityInfoSchema = Yup.object().shape({
    storageFacility: Yup.object().shape({
      facilityName: Yup.string().required('Required'),
      phone: phoneSchema,
      email: emailSchema,
      serviceOrderNumber: Yup.string()
        .required('Required')
        .matches(/^[0-9a-zA-Z]+$/, 'Letters and numbers only'),
    }),
    storageFacilityAddress: Yup.object().shape({
      address: requiredAddressSchema,
      lotNumber: Yup.string(),
    }),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.EditFacilityInfoModal}>
          <ShipmentTag shipmentType={shipmentType} />
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Edit facility info and address</h2>
          </ModalTitle>
          <Formik
            validationSchema={editFacilityInfoSchema}
            onSubmit={onSubmit}
            initialValues={{
              storageFacility,
              storageFacilityAddress,
            }}
          >
            {({ isValid }) => {
              return (
                <Form className={formStyles.form}>
                  <Fieldset>
                    <h3 className={styles.ModalSubTitle}>Facility info</h3>
                    <Grid row>
                      <Grid col={12}>
                        <TextField label="Facility name" id="facilityName" name="storageFacility.facilityName" />
                      </Grid>
                    </Grid>

                    <Grid row gap>
                      <Grid col={6}>
                        <MaskedTextField
                          label="Phone"
                          id="facilityPhone"
                          name="storageFacility.phone"
                          type="tel"
                          minimum="12"
                          mask="000{-}000{-}0000"
                          optional
                        />
                      </Grid>
                    </Grid>

                    <Grid row>
                      <Grid col={12}>
                        <TextField label="Email" id="facilityEmail" name="storageFacility.email" optional />
                      </Grid>
                    </Grid>

                    <Grid row gap>
                      <Grid col={6}>
                        <FormGroup>
                          <TextField
                            label="Service order number"
                            id="facilityServiceOrderNumber"
                            name="storageFacility.serviceOrderNumber"
                          />
                        </FormGroup>
                      </Grid>
                    </Grid>
                  </Fieldset>
                  <Fieldset>
                    <h3 className={styles.ModalSubTitle}>Storage facility address</h3>
                    <AddressFields
                      name="storageFacilityAddress.address"
                      className={styles.AddressFields}
                      render={(fields) => (
                        <>
                          {fields}
                          <Grid row gap>
                            <Grid col={6}>
                              <FormGroup>
                                <Label htmlFor="facilityLotNumber">
                                  Lot number
                                  <span className="float-right usa-hint">Optional</span>
                                </Label>
                                <Field as={TextInput} id="facilityLotNumber" name="storageFacilityAddress.lotNumber" />
                              </FormGroup>
                            </Grid>
                          </Grid>
                        </>
                      )}
                    />
                  </Fieldset>
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

EditFacilityInfoModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  storageFacility: StorageFacilityShape.isRequired,
  storageFacilityAddress: StorageFacilityAddressShape.isRequired,
  shipmentType: ShipmentOptionsOneOf.isRequired,
};
export default EditFacilityInfoModal;
