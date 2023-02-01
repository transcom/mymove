import React, { useEffect, useState } from 'react';
import { useMutation } from 'react-query';
import { func, number, object } from 'prop-types';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Alert, FormGroup, Label, Radio, Textarea } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewWeightTicket.module.scss';

import { ErrorMessage, Form } from 'components/form';
import { patchWeightTicket } from 'services/ghcApi';
import { ShipmentShape, WeightTicketShape } from 'types/shipment';
import Fieldset from 'shared/Fieldset';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import formStyles from 'styles/form.module.scss';
import approveRejectStyles from 'styles/approveRejectControls.module.scss';
import { formatWeight } from 'utils/formatters';
import ppmDocumentStatus from 'constants/ppms';

const validationSchema = Yup.object().shape({
  emptyWeight: Yup.number().required('Enter the empty weight'),
  fullWeight: Yup.number()
    .required('Required')
    .when('emptyWeight', (emptyWeight, schema) => {
      return emptyWeight != null
        ? schema.min(emptyWeight + 1, 'The full weight must be greater than the empty weight')
        : schema;
    }),
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
  weightTicket,
  tripNumber,
  ppmNumber,
  onError,
  onSuccess,
  formRef,
}) {
  const [canEditRejection, setCanEditRejection] = useState(true);

  const [patchWeightTicketMutation] = useMutation(patchWeightTicket, {
    onSuccess,
    onError,
  });

  const ppmShipment = mtoShipment?.ppmShipment;

  const handleSubmit = (values) => {
    const ownsTrailer = values.ownsTrailer === 'true';
    const trailerMeetsCriteria = ownsTrailer ? values.trailerMeetsCriteria === 'true' : false;
    const payload = {
      ppmShipmentId: weightTicket.ppmShipmentId,
      vehicleDescription: weightTicket.vehicleDescription,
      emptyWeight: parseInt(values.emptyWeight, 10),
      missingEmptyWeightTicket: weightTicket.missingEmptyWeightTicket,
      fullWeight: parseInt(values.fullWeight, 10),
      missingFullWeightTicket: weightTicket.missingFullWeightTicket,
      ownsTrailer,
      trailerMeetsCriteria,
      reason: values.status === 'APPROVED' ? null : values.rejectionReason,
      status: values.status,
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

  const initialValues = {
    emptyWeight: emptyWeight ? `${emptyWeight}` : '',
    fullWeight: fullWeight ? `${fullWeight}` : '',
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
            handleChange(event);
            setCanEditRejection(true);
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
                mask={Number}
                description={missingEmptyWeightTicket ? 'Constructed weight' : 'Weight tickets'}
                scale={0} // digits after point, 0 for integers
                signed={false} // disallow negative
                thousandsSeparator=","
                lazy={false} // immediate masking evaluation
                suffix="lbs"
              />

              <MaskedTextField
                defaultValue="0"
                name="fullWeight"
                label="Full weight"
                id="fullWeight"
                mask={Number}
                description={missingFullWeightTicket ? 'Constructed weight' : 'Weight tickets'}
                scale={0} // digits after point, 0 for integers
                signed={false} // disallow negative
                thousandsSeparator=","
                lazy={false} // immediate masking evaluation
                suffix="lbs"
              />
              <Label className={styles.label}>Net weight</Label>
              <div className={styles.displayValue}>{formatWeight(values.fullWeight - values.emptyWeight)}</div>
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
              <p>Add a review for this weight ticket</p>
              <ErrorMessage display={!!errors?.status && !!touched?.status}>{errors.status}</ErrorMessage>
              <Fieldset>
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
                    onChange={handleChange}
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
                            name="rejectionReason"
                            onChange={handleChange}
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
};

ReviewWeightTicket.defaultProps = {
  weightTicket: null,
  mtoShipment: null,
  onSuccess: null,
  formRef: null,
};
