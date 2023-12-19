import React, { useEffect, useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { func, number, object, PropTypes } from 'prop-types';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Alert, FormGroup, Label, Radio, Textarea } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';
import EditPPMNetWeight from '../EditNetWeights/EditPPMNetWeight';

import styles from './ReviewWeightTicket.module.scss';

import { removeCommas } from 'utils/formatters';
import { ErrorMessage, Form } from 'components/form';
import { patchWeightTicket } from 'services/ghcApi';
import { ShipmentShape, WeightTicketShape } from 'types/shipment';
import { OrderShape } from 'types/order';
import Fieldset from 'shared/Fieldset';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import formStyles from 'styles/form.module.scss';
import approveRejectStyles from 'styles/approveRejectControls.module.scss';
import ppmDocumentStatus from 'constants/ppms';
import { getReimbursableWeight } from 'utils/shipmentWeights';
import { calculateWeightRequested } from 'hooks/custom';
import { isNullUndefinedOrWhitespace } from 'shared/utils';

const validationSchema = Yup.object().shape({
  emptyWeight: Yup.number().required('Enter the empty weight'),
  fullWeight: Yup.number()
    .required('Required')
    .when('emptyWeight', ([emptyWeight], schema) => {
      return emptyWeight != null
        ? schema.min(emptyWeight + 1, 'The full weight must be greater than the empty weight')
        : schema;
    }),
  reimbursableWeight: Yup.number().required('Required').min(0, 'reimbursable weight must be at least 0'),
  trailerMeetsCriteria: Yup.string().when('ownsTrailer', {
    is: 'true',
    then: (schema) => schema.required('Required'),
  }),
  rejectionReason: Yup.string().when('status', {
    is: ppmDocumentStatus.REJECTED,
    then: (schema) => schema.required('Add a reason why this weight ticket is rejected'),
  }),
  status: Yup.string().required('Reviewing this weight ticket is required'),
});

export default function ReviewWeightTicket({
  mtoShipment,
  mtoShipments,
  order,
  weightTicket,
  tripNumber,
  ppmNumber,
  onError,
  onSuccess,
  formRef,
  updateTotalWeight,
}) {
  const [canEditRejection, setCanEditRejection] = useState(true);
  const { mutate: patchWeightTicketMutation } = useMutation({
    mutationFn: patchWeightTicket,
    onSuccess,
    onError,
  });
  const ppmShipment = mtoShipment?.ppmShipment;

  const weightAllowance = order.entitlement?.totalWeight;

  const handleSubmit = (formValues) => {
    const ownsTrailer = formValues.ownsTrailer === 'true';
    const trailerMeetsCriteria = ownsTrailer ? formValues.trailerMeetsCriteria === 'true' : false;
    const payload = {
      ppmShipmentId: weightTicket.ppmShipmentId,
      vehicleDescription: weightTicket.vehicleDescription,
      emptyWeight: parseInt(removeCommas(formValues.emptyWeight), 10),
      missingEmptyWeightTicket: weightTicket.missingEmptyWeightTicket,
      fullWeight: parseInt(removeCommas(formValues.fullWeight), 10),
      missingFullWeightTicket: weightTicket.missingFullWeightTicket,
      ownsTrailer,
      trailerMeetsCriteria,
      reason: formValues.rejectionReason,
      status: formValues.status,
      reimbursableWeight: parseInt(removeCommas(formValues.reimbursableWeight), 10),
    };
    patchWeightTicketMutation({
      ppmShipmentId: weightTicket.ppmShipmentId,
      weightTicketId: weightTicket.id,
      payload,
      eTag: weightTicket.eTag,
    });
  };

  const {
    vehicleDescription,
    missingEmptyWeightTicket,
    missingFullWeightTicket,
    emptyWeight,
    fullWeight,
    reimbursableWeight = getReimbursableWeight(weightTicket),
    ownsTrailer,
    proofOfTrailerOwnershipDocument,
    trailerMeetsCriteria,
    status,
    reason,
  } = weightTicket || {};

  const hasProofOfTrailerOwnershipDocument = proofOfTrailerOwnershipDocument?.uploads.length > 0;
  let isTrailerClaimable;
  if (ownsTrailer) {
    isTrailerClaimable = trailerMeetsCriteria ? 'true' : 'false';
  } else {
    isTrailerClaimable = '';
  }

  const updateMtoShipmentsWithNewWeightValues = (currentMtoShipments, values) => {
    const mtoShipmentIndex = currentMtoShipments.findIndex((index) => index.id === mtoShipment.id);
    const updatedWeightTicket = {
      ...weightTicket,
      emptyWeight: parseInt(removeCommas(values.emptyWeight), 10),
      fullWeight: parseInt(removeCommas(values.fullWeight), 10),
      reimbursableWeight: parseInt(removeCommas(values.reimbursableWeight), 10),
    };
    const updatedPPMShipment = {
      ...currentMtoShipments[mtoShipmentIndex].ppmShipment,
    };
    const weightTicketIndex = updatedPPMShipment.weightTickets.findIndex(
      (ticket) => ticket.id === updatedWeightTicket.id,
    );
    updatedPPMShipment.weightTickets[weightTicketIndex] = updatedWeightTicket;
    const updatedMtoShipment = {
      ...mtoShipment,
      ppmShipment: updatedPPMShipment,
    };
    const updatedMtoShipments = currentMtoShipments;
    updatedMtoShipments[mtoShipmentIndex] = updatedMtoShipment;
    return updatedMtoShipments;
  };
  const getNewNetWeightCalculation = (currentMtoShipments, updatedFormValues) => {
    const newMtoShipments = updateMtoShipmentsWithNewWeightValues(currentMtoShipments, updatedFormValues);
    return calculateWeightRequested(newMtoShipments);
  };

  const initialValues = {
    emptyWeight: emptyWeight ? `${emptyWeight}` : '',
    fullWeight: fullWeight ? `${fullWeight}` : '',
    reimbursableWeight: reimbursableWeight ? `${reimbursableWeight}` : '',
    ownsTrailer: ownsTrailer ? 'true' : 'false',
    trailerMeetsCriteria: isTrailerClaimable,
    status: status || '',
    rejectionReason: reason || '',
  };

  useEffect(() => {
    if (formRef?.current) {
      formRef.current.resetForm();
      formRef.current.validateForm();
    }
  }, [formRef, weightTicket]);

  return (
    <div className={classnames(styles.container, 'container--accent--ppm')}>
      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema}
        innerRef={formRef}
        onSubmit={handleSubmit}
        enableReinitialize
        validateOnMount
      >
        {({ handleChange, errors, setFieldError, setFieldTouched, setFieldValue, touched, values }) => {
          const handleApprovalChange = (event) => {
            const updatedValue = event.target.value;
            const newApprovalState =
              updatedValue === ppmDocumentStatus.APPROVED && !isNullUndefinedOrWhitespace(updatedValue);
            if (newApprovalState === false) {
              setCanEditRejection(true);
              setFieldValue('rejectionReason', '');
            }
            handleChange(event);
          };
          const handleRejectionReasonChange = (event) => {
            handleChange(event);
          };
          const handleFieldValueChange = (event) => {
            setFieldValue(event.target.name, removeCommas(event.target.value));
            if (mtoShipments !== undefined && mtoShipments.length > 0) {
              updateTotalWeight(getNewNetWeightCalculation(mtoShipments, values));
            }
          };
          const handleTrailerOwnedChange = (event) => {
            handleChange(event);
            setFieldValue('trailerMeetsCriteria', '');
            setFieldTouched('trailerMeetsCriteria', false, false);
            setFieldError('trailerMeetsCriteria', null);
          };

          const handleTrailerClaimableChange = (event) => {
            if (event.target.value === 'true') {
              setFieldValue('status', '');
            }
            handleChange(event);
          };

          return (
            <Form className={classnames(formStyles.form, styles.ReviewWeightTicket)}>
              <PPMHeaderSummary ppmShipment={ppmShipment} ppmNumber={ppmNumber} />
              <hr />
              <h3 className={styles.tripNumber}>Trip {tripNumber}</h3>
              <legend className={classnames('usa-label', styles.label)}>Vehicle description</legend>
              <div className={styles.displayValue}>{vehicleDescription}</div>

              <MaskedTextField
                defaultValue="0"
                name="emptyWeight"
                label="Empty weight"
                id="emptyWeight"
                inputTestId="emptyWeight"
                mask={Number}
                description={missingEmptyWeightTicket ? 'Vehicle weight' : 'Weight tickets'}
                scale={0} // digits after point, 0 for integers
                signed={false} // disallow negative
                thousandsSeparator=","
                lazy={false} // immediate masking evaluation
                suffix="lbs"
                onBlur={handleFieldValueChange}
              />

              <MaskedTextField
                defaultValue="0"
                name="fullWeight"
                label="Full weight"
                id="fullWeight"
                inputTestId="fullWeight"
                mask={Number}
                description={missingFullWeightTicket ? 'Constructed weight' : 'Weight tickets'}
                scale={0} // digits after point, 0 for integers
                signed={false} // disallow negative
                thousandsSeparator=","
                lazy={false} // immediate masking evaluation
                suffix="lbs"
                onBlur={handleFieldValueChange}
              />

              <MaskedTextField
                defaultValue="0"
                name="reimbursableWeight"
                label="Reimbursable weight"
                id="reimbursableWeight"
                inputTestId="reimbursableWeight"
                mask={Number}
                description="Maximum allowable weight"
                scale={0} // digits after point, 0 for integers
                signed={false} // disallow negative
                thousandsSeparator=","
                lazy={false} // immediate masking evaluation
                suffix="lbs"
                onBlur={handleFieldValueChange}
              />

              <EditPPMNetWeight
                weightTicket={weightTicket}
                weightAllowance={weightAllowance}
                shipments={mtoShipments}
              />

              <FormGroup>
                <Fieldset>
                  <legend className="usa-label">Did they use a trailer they owned?</legend>
                  <Field
                    as={Radio}
                    id="ownsTrailerYes"
                    label="Yes"
                    name="ownsTrailer"
                    value="true"
                    checked={values.ownsTrailer === 'true'}
                    onChange={handleTrailerOwnedChange}
                  />
                  <Field
                    as={Radio}
                    id="ownsTrailerNo"
                    label="No"
                    name="ownsTrailer"
                    value="false"
                    checked={values.ownsTrailer === 'false'}
                    onChange={handleTrailerOwnedChange}
                  />
                </Fieldset>
              </FormGroup>
              {values.ownsTrailer === 'true' && (
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">{`Is the trailer's weight claimable?`}</legend>
                    <ErrorMessage display={!!errors?.trailerMeetsCriteria && !!touched?.trailerMeetsCriteria}>
                      {errors.trailerMeetsCriteria}
                    </ErrorMessage>
                    <Field
                      as={Radio}
                      id="trailerCriteriaYes"
                      label="Yes"
                      name="trailerMeetsCriteria"
                      value="true"
                      checked={values.trailerMeetsCriteria === 'true'}
                      onChange={handleTrailerClaimableChange}
                    />
                    <Field
                      as={Radio}
                      id="trailerCriteriaNo"
                      label="No"
                      name="trailerMeetsCriteria"
                      value="false"
                      checked={values.trailerMeetsCriteria === 'false'}
                      onChange={handleTrailerClaimableChange}
                    />
                    {values.trailerMeetsCriteria === 'true' && !hasProofOfTrailerOwnershipDocument && (
                      <Alert type="info">Proof of ownership is needed to accept this item.</Alert>
                    )}
                  </Fieldset>
                </FormGroup>
              )}
              <h3 className={styles.reviewHeader}>Review trip {tripNumber}</h3>
              <p>Add a review for this Weight Ticket</p>
              <ErrorMessage display={!!errors?.status && !!touched?.status}>{errors.status}</ErrorMessage>
              <Fieldset className={styles.statusOptions}>
                <div
                  className={classnames(approveRejectStyles.statusOption, {
                    [approveRejectStyles.selected]: values.status === ppmDocumentStatus.APPROVED,
                  })}
                >
                  <Radio
                    id={`approve-${weightTicket?.id}`}
                    checked={values.status === ppmDocumentStatus.APPROVED}
                    value={ppmDocumentStatus.APPROVED}
                    name="status"
                    label="Accept"
                    onChange={handleApprovalChange}
                    data-testid="approveRadio"
                    disabled={values.trailerMeetsCriteria === 'true' && !hasProofOfTrailerOwnershipDocument}
                  />
                </div>
                <div
                  className={classnames(approveRejectStyles.statusOption, styles.reject, {
                    [approveRejectStyles.selected]: values.status === ppmDocumentStatus.REJECTED,
                  })}
                >
                  <Radio
                    id={`reject-${weightTicket?.id}`}
                    checked={values.status === ppmDocumentStatus.REJECTED}
                    value={ppmDocumentStatus.REJECTED}
                    name="status"
                    label="Reject"
                    onChange={handleApprovalChange}
                    data-testid="rejectRadio"
                  />

                  {values.status === ppmDocumentStatus.REJECTED && (
                    <FormGroup className={styles.reason}>
                      <Label htmlFor={`rejectReason-${weightTicket?.id}`}>Reason</Label>
                      {!canEditRejection && (
                        <p data-testid="rejectionReasonReadOnly">{weightTicket?.reason || values.rejectionReason}</p>
                      )}

                      {canEditRejection && (
                        <>
                          <ErrorMessage display={!!errors?.rejectionReason && !!touched?.rejectionReason}>
                            {errors.rejectionReason}
                          </ErrorMessage>
                          <Textarea
                            id={`rejectReason-${weightTicket?.id}`}
                            data-testid="rejectionReasonText"
                            name="rejectionReason"
                            onChange={handleRejectionReasonChange}
                            error={touched.rejectionReason ? errors.rejectionReason : null}
                            value={values.rejectionReason}
                            placeholder="Type something"
                          />
                          <div className={styles.hint}>{500 - values.rejectionReason.length} characters</div>
                        </>
                      )}
                    </FormGroup>
                  )}
                </div>
              </Fieldset>
            </Form>
          );
        }}
      </Formik>
    </div>
  );
}

ReviewWeightTicket.propTypes = {
  weightTicket: WeightTicketShape,
  mtoShipment: ShipmentShape,
  tripNumber: number.isRequired,
  ppmNumber: number.isRequired,
  onSuccess: func,
  formRef: object,
  mtoShipments: PropTypes.arrayOf(ShipmentShape),
  order: OrderShape.isRequired,
};

ReviewWeightTicket.defaultProps = {
  weightTicket: null,
  mtoShipment: null,
  onSuccess: null,
  formRef: null,
  mtoShipments: [],
};
