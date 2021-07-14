import React from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Radio, FormGroup, Label } from '@trussworks/react-uswds';

import { Form } from 'components/form/Form';
import FileUpload from 'components/FileUpload/FileUpload';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { documentSizeLimitMsg } from 'shared/constants';
import profileImage from 'scenes/Review/images/profile.png';
import Hint from 'components/Hint/index';
import { DropdownArrayOf, ExistingUploadsShape } from 'types';
import { DutyStationShape } from 'types/dutyStation';
import { OrdersShape } from 'types/customerShapes';
import { DropdownInput, DatePickerInput, DutyStationInput } from 'components/form/fields';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import formStyles from 'styles/form.module.scss';

const EditOrdersForm = ({
  createUpload,
  onDelete,
  initialValues,
  existingUploads,
  onUploadComplete,
  filePondEl,
  onSubmit,
  ordersTypeOptions,
  currentStation,
  onCancel,
}) => {
  const validationSchema = Yup.object().shape({
    orders_type: Yup.mixed()
      .oneOf(ordersTypeOptions.map((i) => i.key))
      .required('Required'),
    issue_date: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),
    report_by_date: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),
    has_dependents: Yup.mixed().oneOf(['yes', 'no']).required('Required'),
    new_duty_station: Yup.object()
      .shape({
        name: Yup.string().notOneOf(
          [currentStation?.name],
          'You entered the same duty station for your origin and destination. Please change one of them.',
        ),
      })
      .nullable()
      .required('Required'),
  });

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={validationSchema} validateOnMount>
      {({ isValid, isSubmitting, handleSubmit }) => {
        return (
          <Form className={formStyles.form}>
            <img src={profileImage} alt="" />
            <h1
              style={{
                display: 'inline-block',
                marginLeft: 10,
                marginBottom: 16,
                marginTop: 20,
              }}
            >
              Orders
            </h1>
            <SectionWrapper className={formStyles.formSection}>
              <DropdownInput label="Orders type" name="orders_type" options={ordersTypeOptions} required />
              <DatePickerInput
                name="issue_date"
                label="Orders date"
                required
                renderInput={(input) => (
                  <>
                    {input}
                    <Hint>
                      <p>Date your orders were issued.</p>
                    </Hint>
                  </>
                )}
              />
              <DatePickerInput name="report_by_date" label="Report-by date" required />
              <FormGroup>
                <Label>Are dependents included in your orders?</Label>
                <div>
                  <Field
                    as={Radio}
                    label="Yes"
                    id="hasDependentsYes"
                    name="has_dependents"
                    value="yes"
                    title="Yes, dependents are included in my orders"
                    type="radio"
                  />
                  <Field
                    as={Radio}
                    label="No"
                    id="hasDependentsNo"
                    name="has_dependents"
                    value="no"
                    title="No, dependents are not included in my orders"
                    type="radio"
                  />
                </div>
              </FormGroup>
              <DutyStationInput name="new_duty_station" label="New duty station" displayAddress={false} />
              <p>Uploads:</p>
              {existingUploads?.length > 0 && <UploadsTable uploads={existingUploads} onDelete={onDelete} />}
              {initialValues?.uploaded_orders && (
                <div>
                  <p>{documentSizeLimitMsg}</p>
                  <FileUpload
                    ref={filePondEl}
                    createUpload={createUpload}
                    onChange={onUploadComplete}
                    labelIdle={'Drag & drop or <span class="filepond--label-action">click to upload orders</span>'}
                  />
                </div>
              )}
            </SectionWrapper>

            <div className={formStyles.formActions}>
              <WizardNavigation
                editMode
                onCancelClick={onCancel}
                disableNext={!isValid || isSubmitting}
                onNextClick={handleSubmit}
              />
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

EditOrdersForm.propTypes = {
  ordersTypeOptions: DropdownArrayOf.isRequired,
  createUpload: PropTypes.func.isRequired,
  onUploadComplete: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  existingUploads: ExistingUploadsShape,
  filePondEl: PropTypes.shape({
    current: PropTypes.shape({}),
  }),
  initialValues: OrdersShape.isRequired,
  currentStation: DutyStationShape.isRequired,
  onCancel: PropTypes.func.isRequired,
};

EditOrdersForm.defaultProps = {
  existingUploads: [],
  filePondEl: null,
};

export default EditOrdersForm;
