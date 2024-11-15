import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Radio, FormGroup, Label, Link as USWDSLink } from '@trussworks/react-uswds';

import { isBooleanFlagEnabled } from '../../../utils/featureFlags';

import styles from './EditOrdersForm.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import ToolTip from 'shared/ToolTip/ToolTip';
import { ORDERS_PAY_GRADE_OPTIONS, ORDERS_TYPE } from 'constants/orders';
import { Form } from 'components/form/Form';
import FileUpload from 'components/FileUpload/FileUpload';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { FEATURE_FLAG_KEYS, documentSizeLimitMsg } from 'shared/constants';
import profileImage from 'scenes/Review/images/profile.png';
import { DropdownArrayOf } from 'types';
import { ExistingUploadsShape } from 'types/uploads';
import { DropdownInput, DatePickerInput, DutyLocationInput } from 'components/form/fields';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import Callout from 'components/Callout';
import { formatLabelReportByDate, dropdownInputOptions, formatYesNoAPIValue } from 'utils/formatters';
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
  const [currentDutyLocation, setDutyLocation] = useState(initialValues.origin_duty_location);
  const [newDutyLocation, setNewDutyLocation] = useState(initialValues.new_duty_location);
  const [showAccompaniedTourField, setShowAccompaniedTourField] = useState(false);
  const [showDependentAgeFields, setShowDependentAgeFields] = useState(false);
  const [hasDependents, setHasDependents] = useState(formatYesNoAPIValue(initialValues.has_dependents));
  const [isOconusMove, setIsOconusMove] = useState(false);
  const [enableUB, setEnableUB] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [finishedFetchingFF, setFinishedFetchingFF] = useState(false);

  const isInitialHasDependentsDisabled =
    initialValues.orders_type === ORDERS_TYPE.STUDENT_TRAVEL ||
    initialValues.orders_type === ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS;
  const [isHasDependentsDisabled, setHasDependentsDisabled] = useState(isInitialHasDependentsDisabled);
  const [prevOrderType, setPrevOrderType] = useState(initialValues.orders_type);
  const [filteredOrderTypeOptions, setFilteredOrderTypeOptions] = useState(ordersTypeOptions);

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
    counseling_office_id: currentDutyLocation?.provides_services_counseling
      ? Yup.string().required('Required')
      : Yup.string().notRequired(),
    accompanied_tour: showAccompaniedTourField
      ? Yup.mixed().oneOf(['yes', 'no']).required('Required')
      : Yup.string().notRequired(),
    dependents_under_twelve: showDependentAgeFields
      ? Yup.number().min(0).required('Required')
      : Yup.number().notRequired(),
    dependents_twelve_and_over: showDependentAgeFields
      ? Yup.number().min(0).required('Required')
      : Yup.number().notRequired(),
  });

  const enableDelete = () => {
    const isValuePresent = initialValues.move_status === 'DRAFT';
    return isValuePresent;
  };

  const payGradeOptions = dropdownInputOptions(ORDERS_PAY_GRADE_OPTIONS);

  let originMeta;
  let newDutyMeta = '';

  useEffect(() => {
    // Only check the FF on load
    const checkUBFeatureFlag = async () => {
      const enabled = await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.UNACCOMPANIED_BAGGAGE);
      if (enabled) {
        setEnableUB(true);
      }
      setFinishedFetchingFF(true);
    };
    checkUBFeatureFlag();
  }, []);
  useEffect(() => {
    showCounselingOffices(currentDutyLocation?.id).then((fetchedData) => {
      if (fetchedData.body) {
        const counselingOffices = fetchedData.body.map((item) => ({
          key: item.id,
          value: item.name,
        }));
        setOfficeOptions(counselingOffices);
      }
    });
    // Check if either currentDutyLocation or newDutyLocation is OCONUS
    if (currentDutyLocation?.address?.isOconus || newDutyLocation?.address?.isOconus) {
      setIsOconusMove(true);
    } else {
      setIsOconusMove(false);
    }
    if (currentDutyLocation?.address && newDutyLocation?.address && enableUB) {
      // Only if one of the duty locations is OCONUS should accompanied tour and dependent
      // age fields display
      if (isOconusMove && hasDependents) {
        setShowAccompaniedTourField(true);
        setShowDependentAgeFields(true);
      } else {
        setShowAccompaniedTourField(false);
        setShowDependentAgeFields(false);
      }
    }
    if (isLoading && finishedFetchingFF) {
      // If the form is still loading and the FF has finished fetching,
      // then the form is done loading
      setIsLoading(false);
    }
  }, [currentDutyLocation, newDutyLocation, isOconusMove, hasDependents, enableUB, finishedFetchingFF, isLoading]);

  useEffect(() => {
    const fetchData = async () => {
      const alaskaEnabled = await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.ENABLE_ALASKA);

      const updatedOptions = alaskaEnabled
        ? ordersTypeOptions
        : ordersTypeOptions.filter(
            (e) => e.key !== ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS && e.key !== ORDERS_TYPE.STUDENT_TRAVEL,
          );

      setFilteredOrderTypeOptions(updatedOptions);
    };
    fetchData();
  }, [ordersTypeOptions]);

  if (isLoading) {
    return <LoadingPlaceholder />;
  }

  return (
    <Formik
      initialValues={{
        ...initialValues,
        has_dependents: isInitialHasDependentsDisabled ? 'yes' : initialValues.has_dependents,
      }}
      onSubmit={onSubmit}
      validationSchema={validationSchema}
      validateOnMount
      initialTouched={{
        orders_type: true,
        issue_date: true,
        report_by_date: true,
        has_dependents: true,
        grade: true,
        accompanied_tour: true,
        dependents_under_twelve: true,
        dependents_twelve_and_over: true,
        origin_duty_location: true,
        new_duty_location: true,
      }}
    >
      {({ isValid, isSubmitting, handleSubmit, handleChange, values, setFieldValue }) => {
        const isRetirementOrSeparation = ['RETIREMENT', 'SEPARATION'].includes(values.orders_type);

        if (!values.origin_duty_location) originMeta = 'Required';
        else originMeta = null;

        if (!values.new_duty_location) newDutyMeta = 'Required';
        else newDutyMeta = null;

        const handleHasDependentsChange = (e) => {
          // Declare a duplicate local scope of the field value
          // for the form to prevent state race conditions
          if (e.target.value === '') {
            setFieldValue('has_dependents', '');
          } else {
            const fieldValueHasDependents = e.target.value === 'yes';
            setHasDependents(fieldValueHasDependents);
            setFieldValue('has_dependents', fieldValueHasDependents ? 'yes' : 'no');
            if (fieldValueHasDependents && isOconusMove && enableUB) {
              setShowAccompaniedTourField(true);
              setShowDependentAgeFields(true);
            } else {
              setShowAccompaniedTourField(false);
              setShowDependentAgeFields(false);
            }
          }
        };

        const handleOrderTypeChange = (e) => {
          const { value } = e.target;
          if (value === ORDERS_TYPE.STUDENT_TRAVEL || value === ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS) {
            setHasDependentsDisabled(true);
            handleHasDependentsChange({ target: { value: 'yes' } });
          } else {
            setHasDependentsDisabled(false);
            if (
              prevOrderType === ORDERS_TYPE.STUDENT_TRAVEL ||
              prevOrderType === ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS
            ) {
              handleHasDependentsChange({ target: { value: '' } });
            }
          }
          setPrevOrderType(value);
        };

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
                options={filteredOrderTypeOptions}
                required
                hint="Required"
                onChange={(e) => {
                  handleChange(e);
                  handleOrderTypeChange(e);
                }}
              />
              <DatePickerInput name="issue_date" label="Orders date" hint="Required" required />
              <DatePickerInput
                name="report_by_date"
                label={formatLabelReportByDate(values.orders_type)}
                required
                hint="Required"
              />
              <DutyLocationInput
                label="Current duty location"
                name="origin_duty_location"
                id="origin_duty_location"
                hint="Required"
                onDutyLocationChange={(e) => {
                  setDutyLocation(e);
                }}
                required
                metaOverride={originMeta}
              />
              {currentDutyLocation?.provides_services_counseling && (
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
                    onDutyLocationChange={(e) => {
                      setNewDutyLocation(e);
                    }}
                  />
                </>
              ) : (
                <DutyLocationInput
                  name="new_duty_location"
                  label="New duty location"
                  displayAddress={false}
                  metaOverride={newDutyMeta}
                  onDutyLocationChange={(e) => {
                    setNewDutyLocation(e);
                  }}
                />
              )}

              <FormGroup>
                <Label hint="Required">Are dependents included in your orders?</Label>
                <div>
                  <Field
                    as={Radio}
                    label="Yes"
                    id="hasDependentsYes"
                    data-testid="hasDependentsYes"
                    name="has_dependents"
                    value="yes"
                    title="Yes, dependents are included in my orders"
                    type="radio"
                    onChange={(e) => {
                      handleHasDependentsChange(e);
                    }}
                    disabled={isHasDependentsDisabled}
                  />
                  <Field
                    as={Radio}
                    label="No"
                    id="hasDependentsNo"
                    data-testid="hasDependentsNo"
                    name="has_dependents"
                    value="no"
                    title="No, dependents are not included in my orders"
                    type="radio"
                    onChange={(e) => {
                      handleHasDependentsChange(e);
                    }}
                    disabled={isHasDependentsDisabled}
                  />
                </div>
              </FormGroup>

              {showAccompaniedTourField && (
                <FormGroup>
                  <Label hint="Required">Is this an accompanied tour?</Label>
                  <div>
                    <div className={styles.radioWithToolTip}>
                      <Field
                        as={Radio}
                        label="Yes"
                        id="isAnAccompaniedTourYes"
                        data-testid="isAnAccompaniedTourYes"
                        name="accompanied_tour"
                        value="yes"
                        type="radio"
                      />
                      <ToolTip
                        text="Accompanied Tour: An authorized order (assignment or tour) that allows dependents to travel to the new Permanent Duty Station (PDS)"
                        position="right"
                        icon="info-circle"
                        color="blue"
                        data-testid="isAnAccompaniedTourYesToolTip"
                        closeOnLeave
                      />
                    </div>
                    <div className={styles.radioWithToolTip}>
                      <Field
                        as={Radio}
                        label="No"
                        id="isAnAccompaniedTourNo"
                        data-testid="isAnAccompaniedTourNo"
                        name="accompanied_tour"
                        value="no"
                        type="radio"
                      />
                      <ToolTip
                        text="Unaccompanied Tour: An authorized order (assignment or tour) that DOES NOT allow dependents to travel to the new Permanent Duty Station (PDS)"
                        position="right"
                        icon="info-circle"
                        color="blue"
                        data-testid="isAnAccompaniedTourNoToolTip"
                        closeOnLeave
                      />
                    </div>
                  </div>
                </FormGroup>
              )}

              {showDependentAgeFields && (
                <FormGroup>
                  <MaskedTextField
                    data-testid="dependentsUnderTwelve"
                    defaultValue="0"
                    name="dependents_under_twelve"
                    label="Number of dependents under the age of 12"
                    id="dependentsUnderTwelve"
                    labelHint="Required"
                    mask={Number}
                    scale={0}
                    signed={false}
                    thousandsSeparator=","
                    lazy={false}
                  />

                  <MaskedTextField
                    data-testid="dependentsTwelveAndOver"
                    defaultValue="0"
                    name="dependents_twelve_and_over"
                    label="Number of dependents of the age 12 or over"
                    id="dependentsTwelveAndOver"
                    mask={Number}
                    scale={0}
                    signed={false}
                    labelHint="Required"
                    thousandsSeparator=","
                    lazy={false}
                  />
                </FormGroup>
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
    dependents_under_twelve: PropTypes.string,
    dependents_twelve_and_over: PropTypes.string,
    accompanied_tour: PropTypes.string,
    uploaded_orders: ExistingUploadsShape,
  }).isRequired,
  onCancel: PropTypes.func.isRequired,
};

EditOrdersForm.defaultProps = {
  filePondEl: null,
};

export default EditOrdersForm;
