import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Radio, FormGroup, Label, Link as USWDSLink } from '@trussworks/react-uswds';

import styles from './EditOrdersForm.module.scss';

import { ORDERS_PAY_GRADE_OPTIONS } from 'constants/orders';
import { Form } from 'components/form/Form';
import FileUpload from 'components/FileUpload/FileUpload';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { documentSizeLimitMsg } from 'shared/constants';
import profileImage from 'scenes/Review/images/profile.png';
import { DropdownArrayOf } from 'types';
import { ExistingUploadsShape } from 'types/uploads';
import { DropdownInput, DatePickerInput, DutyLocationInput } from 'components/form/fields';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import Callout from 'components/Callout';
import { formatLabelReportByDate, dropdownInputOptions } from 'utils/formatters';
import formStyles from 'styles/form.module.scss';
import { showCounselingOffices } from 'services/internalApi';

const EditOrdersForm = ({
  createUpload,
  onDelete,
  initialValues,
  onUploadComplete,
  filePondEl,
  onSubmit,
  ordersTypeOptions,
  onCancel,
}) => {
  const [officeOptions, setOfficeOptions] = useState(null);
  const [dutyLocation, setDutyLocation] = useState(initialValues.origin_duty_location);
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
    new_duty_location: Yup.object().nullable().required('Required'),
    uploaded_orders: Yup.array()
      .of(
        Yup.object().shape({
          id: Yup.string(),
          created_at: Yup.string(),
          bytes: Yup.string(),
          url: Yup.string(),
          filename: Yup.string(),
        }),
      )
      .min(1),
    grade: Yup.mixed().oneOf(Object.keys(ORDERS_PAY_GRADE_OPTIONS)).required('Required'),
    origin_duty_location: Yup.object().nullable().required('Required'),
    counseling_office_id: dutyLocation?.provides_services_counseling
      ? Yup.string().required('Required')
      : Yup.string().notRequired(),
  });

  const enableDelete = () => {
    const isValuePresent = initialValues.move_status === 'DRAFT';
    return isValuePresent;
  };

  const payGradeOptions = dropdownInputOptions(ORDERS_PAY_GRADE_OPTIONS);

  let originMeta;
  let newDutyMeta = '';

  useEffect(() => {
    showCounselingOffices(dutyLocation?.id).then((fetchedData) => {
      if (fetchedData.body) {
        const counselingOffices = fetchedData.body.map((item) => ({
          key: item.id,
          value: item.name,
        }));
        setOfficeOptions(counselingOffices);
      }
    });
  }, [dutyLocation]);

  return (
    <Formik
      initialValues={initialValues}
      onSubmit={onSubmit}
      validationSchema={validationSchema}
      validateOnMount
      initialTouched={{ orders_type: true, issue_date: true, report_by_date: true, has_dependents: true, grade: true }}
    >
      {({ isValid, isSubmitting, handleSubmit, setValues, values }) => {
        const isRetirementOrSeparation = ['RETIREMENT', 'SEPARATION'].includes(values.orders_type);

        const handleCounselingOfficeChange = () => {
          setValues({
            ...values,
            counseling_office_id: null,
          });
          setOfficeOptions(null);
        };
        if (!values.origin_duty_location) originMeta = 'Required';
        else originMeta = null;

        if (!values.new_duty_location) newDutyMeta = 'Required';
        else newDutyMeta = null;

        return (
          <Form className={`${formStyles.form} ${styles.EditOrdersForm}`}>
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
              <h2>Edit Orders:</h2>
              <DropdownInput
                label="Orders type"
                name="orders_type"
                options={ordersTypeOptions}
                required
                hint="Required"
              />
              <DatePickerInput name="issue_date" label="Orders date" hint="Required" required />
              <DatePickerInput
                name="report_by_date"
                label={formatLabelReportByDate(values.orders_type)}
                required
                hint="Required"
              />
              <FormGroup>
                <Label hint="Required">Are dependents included in your orders?</Label>
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

              <DutyLocationInput
                label="Current duty location"
                name="origin_duty_location"
                id="origin_duty_location"
                hint="Required"
                onDutyLocationChange={(e) => {
                  setDutyLocation(e);
                  handleCounselingOfficeChange();
                }}
                required
                metaOverride={originMeta}
              />
              {dutyLocation?.provides_services_counseling && (
                <div>
                  <Label>
                    Select an origin duty location that most closely represents your current physical location, not
                    where your shipment will originate, if different. This will allow a nearby transportation office to
                    assist
                  </Label>
                  <DropdownInput
                    label="Counseling Office"
                    name="counseling_office_id"
                    id="counseling_office_id"
                    hint="Required"
                    required
                    options={officeOptions}
                  />
                </div>
              )}
              {isRetirementOrSeparation ? (
                <>
                  <h3 className={styles.calloutLabel}>Where are you entitled to move?</h3>
                  <Callout>
                    <span>The government will pay for your move to:</span>
                    <ul>
                      <li>Home of record (HOR)</li>
                      <li>Place entered active duty (PLEAD)</li>
                    </ul>
                    <p>
                      It might pay for a move to your Home of selection (HOS), anywhere in CONUS. Check your orders.
                    </p>
                    <p>
                      Read more about where you are entitled to move when leaving the military on{' '}
                      <USWDSLink
                        target="_blank"
                        rel="noopener noreferrer"
                        href="https://www.militaryonesource.mil/military-life-cycle/separation-transition/military-separation-retirement/deciding-where-to-live-when-you-leave-the-military/"
                      >
                        Military OneSource.
                      </USWDSLink>
                    </p>
                  </Callout>
                  <DutyLocationInput
                    name="new_duty_location"
                    label="HOR, PLEAD or HOS"
                    displayAddress={false}
                    hint="Enter the option closest to your destination. Your move counselor will identify if there might be a cost to you."
                    placeholder="Enter a city or ZIP"
                    metaOverride={newDutyMeta}
                  />
                </>
              ) : (
                <DutyLocationInput
                  name="new_duty_location"
                  label="New duty location"
                  displayAddress={false}
                  metaOverride={newDutyMeta}
                />
              )}
              <DropdownInput
                label="Pay grade"
                name="grade"
                id="grade"
                required
                options={payGradeOptions}
                hint="Required"
              />

              <p>Uploads:</p>
              <UploadsTable
                uploads={initialValues.uploaded_orders}
                onDelete={onDelete}
                showDeleteButton={enableDelete(initialValues)}
                showDownloadLink
              />
              <div>
                <p>{documentSizeLimitMsg}</p>
                <FileUpload
                  ref={filePondEl}
                  createUpload={createUpload}
                  onChange={onUploadComplete}
                  labelIdle={'Drag & drop or <span class="filepond--label-action">click to upload orders</span>'}
                />
              </div>
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
  filePondEl: PropTypes.shape({
    current: PropTypes.shape({}),
  }),
  initialValues: PropTypes.shape({
    orders_type: PropTypes.string,
    issue_date: PropTypes.string,
    report_by_date: PropTypes.string,
    has_dependents: PropTypes.string,
    new_duty_location: PropTypes.shape({
      name: PropTypes.string,
    }),
    origin_duty_location: PropTypes.shape({
      name: PropTypes.string,
    }),
    counseling_office_id: PropTypes.string,
    uploaded_orders: ExistingUploadsShape,
  }).isRequired,
  onCancel: PropTypes.func.isRequired,
};

EditOrdersForm.defaultProps = {
  filePondEl: null,
};

export default EditOrdersForm;
