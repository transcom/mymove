import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Radio, FormGroup, Label, Link as USWDSLink } from '@trussworks/react-uswds';
import { connect } from 'react-redux';

import { isBooleanFlagEnabled } from '../../../utils/featureFlags';

import styles from './EditOrdersForm.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import ToolTip from 'shared/ToolTip/ToolTip';
import { ORDERS_PAY_GRADE_TYPE, ORDERS_TYPE } from 'constants/orders';
import {
  civilianTDYUBAllowanceWeightWarning,
  FEATURE_FLAG_KEYS,
  MOVE_STATUSES,
  documentSizeLimitMsg,
} from 'shared/constants';
import { Form } from 'components/form/Form';
import FileUpload from 'components/FileUpload/FileUpload';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import profileImage from 'scenes/Review/images/profile.png';
import { DropdownArrayOf } from 'types';
import { ExistingUploadsShape } from 'types/uploads';
import { DropdownInput, DatePickerInput, DutyLocationInput } from 'components/form/fields';
import RequiredAsterisk, { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import Callout from 'components/Callout';
import { formatLabelReportByDate, formatPayGradeOptions, formatYesNoAPIValue } from 'utils/formatters';
import formStyles from 'styles/form.module.scss';
import { getPayGradeOptions, getRankOptions, showCounselingOffices } from 'services/internalApi';
import { setShowLoadingSpinner as setShowLoadingSpinnerAction } from 'store/general/actions';
import { milmoveLogger } from 'utils/milmoveLog';
import retryPageLoading from 'utils/retryPageLoading';
import Hint from 'components/Hint';
import { sortRankOptions } from 'shared/utils';
import { selectServiceMemberAffiliation } from 'store/entities/selectors';

const EditOrdersForm = ({
  createUpload,
  onDelete,
  initialValues,
  onUploadComplete,
  filePondEl,
  onSubmit,
  ordersTypeOptions,
  onCancel,
  setShowLoadingSpinner,
  affiliation,
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
  const [ordersType, setOrdersType] = useState(initialValues.orders_type);
  const [grade, setGrade] = useState(initialValues.grade);
  const [isCivilianTDYMove, setIsCivilianTDYMove] = useState(false);
  const [showCivilianTDYUBTooltip, setShowCivilianTDYUBTooltip] = useState(false);
  const [rankOptions, setRankOptions] = useState([]);

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
    grade: Yup.string().required('Required'),
    rank: Yup.mixed()
      .oneOf(rankOptions.map((i) => i.key))
      .required('Required'),
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
    civilian_tdy_ub_allowance: isCivilianTDYMove
      ? Yup.number()
          .transform((value) => (Number.isNaN(value) ? 0 : value))
          .min(0, 'UB weight allowance must be 0 or more')
          .max(2000, 'UB weight allowance cannot exceed 2,000 lbs.')
      : Yup.number().notRequired(),
  });

  const enableDelete = () => {
    const isValuePresent = initialValues.move_status === MOVE_STATUSES.DRAFT;
    return isValuePresent;
  };

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
    const fetchRankOptions = async () => {
      setShowLoadingSpinner(true, null);
      try {
        const fetchedRanks = await getRankOptions(affiliation, grade);
        if (fetchedRanks) {
          const formattedOptions = sortRankOptions(fetchedRanks);
          setRankOptions(formattedOptions);
        }
      } catch (error) {
        const { message } = error;
        milmoveLogger.error({ message, info: null });
        retryPageLoading(error);
      }
      setShowLoadingSpinner(false, null);
    };

    fetchRankOptions();
  }, [affiliation, setShowLoadingSpinner, grade]);

  useEffect(() => {
    const fetchCounselingOffices = async () => {
      if (currentDutyLocation?.id && !officeOptions) {
        setShowLoadingSpinner(true, null);
        try {
          const fetchedData = await showCounselingOffices(currentDutyLocation.id);
          if (fetchedData.body) {
            const counselingOffices = fetchedData.body.map((item) => ({
              key: item.id,
              value: item.name,
            }));
            setOfficeOptions(counselingOffices);
          }
        } catch (error) {
          const { message } = error;
          milmoveLogger.error({ message, info: null });
          retryPageLoading(error);
        }
        setShowLoadingSpinner(false, null);
      }
    };
    fetchCounselingOffices();
  }, [currentDutyLocation.id, officeOptions, setShowLoadingSpinner]);

  const [payGradeOptions, setPayGradeOptions] = useState([]);
  useEffect(() => {
    const fetchGradeOptions = async () => {
      setShowLoadingSpinner(true, null);
      try {
        const fetchedRanks = await getPayGradeOptions(affiliation);
        if (fetchedRanks) {
          setPayGradeOptions(formatPayGradeOptions(fetchedRanks.body));
        }
      } catch (error) {
        const { message } = error;
        milmoveLogger.error({ message, info: null });
        retryPageLoading(error);
      }
      setShowLoadingSpinner(false, null);
    };

    fetchGradeOptions();
  }, [affiliation, setShowLoadingSpinner]);

  useEffect(() => {
    // Check if either currentDutyLocation or newDutyLocation is OCONUS
    if (currentDutyLocation?.address?.isOconus || newDutyLocation?.address?.isOconus) {
      setIsOconusMove(true);
    } else {
      setIsOconusMove(false);
    }

    if (currentDutyLocation?.address && newDutyLocation?.address && enableUB) {
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
  }, [currentDutyLocation, newDutyLocation, isOconusMove, hasDependents, enableUB, isLoading, finishedFetchingFF]);

  useEffect(() => {
    if (ordersType && grade && currentDutyLocation?.address && newDutyLocation?.address && enableUB) {
      if (
        isOconusMove &&
        ordersType === ORDERS_TYPE.TEMPORARY_DUTY &&
        grade === ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE
      ) {
        setIsCivilianTDYMove(true);
      } else {
        setIsCivilianTDYMove(false);
      }
    }
  }, [
    currentDutyLocation,
    newDutyLocation,
    isOconusMove,
    hasDependents,
    enableUB,
    ordersType,
    grade,
    isCivilianTDYMove,
  ]);
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
        civilian_tdy_ub_allowance: true,
      }}
    >
      {({ isValid, isSubmitting, handleSubmit, handleChange, setValues, values, setFieldValue }) => {
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
          setOrdersType(value);
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

        // Conditionally set the civilian TDY UB allowance warning message based on provided weight being in the 351 to 2000 lb range
        const showcivilianTDYUBAllowanceWarning =
          values.civilian_tdy_ub_allowance > 350 && values.civilian_tdy_ub_allowance <= 2000;

        let civilianTDYUBAllowanceWarning = '';
        if (showcivilianTDYUBAllowanceWarning) {
          civilianTDYUBAllowanceWarning = civilianTDYUBAllowanceWeightWarning;
        }

        const toggleCivilianTDYUBTooltip = () => {
          setShowCivilianTDYUBTooltip((prev) => !prev);
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
              {requiredAsteriskMessage}
              <DropdownInput
                label="Orders type"
                name="orders_type"
                options={ordersTypeOptions}
                required
                showRequiredAsterisk
                onChange={(e) => {
                  handleChange(e);
                  handleOrderTypeChange(e);
                }}
              />
              <DatePickerInput name="issue_date" label="Orders date" showRequiredAsterisk required />
              <DatePickerInput
                name="report_by_date"
                label={formatLabelReportByDate(values.orders_type)}
                required
                showRequiredAsterisk
              />
              <DutyLocationInput
                label="Current duty location"
                name="origin_duty_location"
                id="origin_duty_location"
                showRequiredAsterisk
                onDutyLocationChange={(e) => {
                  setDutyLocation(e);
                  handleCounselingOfficeChange();
                }}
                required
                metaOverride={originMeta}
              />
              {currentDutyLocation?.provides_services_counseling && (
                <div>
                  <DropdownInput
                    label="Counseling office"
                    name="counseling_office_id"
                    id="counseling_office_id"
                    showRequiredAsterisk
                    required
                    options={officeOptions}
                  />
                  <Hint>
                    Select an origin duty location that most closely represents your current physical location, not
                    where your shipment will originate, if different. This will allow a nearby transportation office to
                    assist.
                  </Hint>
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
                    label="Destination Location (As Authorized on Orders)"
                    displayAddress={false}
                    showRequiredAsterisk
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
                  showRequiredAsterisk
                  displayAddress={false}
                  metaOverride={newDutyMeta}
                  onDutyLocationChange={(e) => {
                    setNewDutyLocation(e);
                  }}
                />
              )}

              <FormGroup>
                <Label>
                  <span>
                    Are dependents included in your orders? <RequiredAsterisk />
                  </span>
                </Label>
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
                  <Label>
                    <span>
                      Is this an accompanied tour? <RequiredAsterisk />
                    </span>
                  </Label>
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
                    showRequiredAsterisk
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
                    showRequiredAsterisk
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
                showRequiredAsterisk
                onChange={(e) => {
                  setGrade(e.target.value);
                  handleChange(e);
                  setFieldValue('rank', '');
                }}
              />

              {grade !== '' ? (
                <DropdownInput
                  label="Rank"
                  name="rank"
                  id="rank"
                  required
                  options={rankOptions}
                  showRequiredAsterisk
                  onChange={(e) => {
                    handleChange(e);
                  }}
                />
              ) : null}

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
                  labelIdle='Drag & drop or <span class="filepond--label-action">click to upload orders</span>'
                />
              </div>

              {isCivilianTDYMove && (
                <FormGroup>
                  <div>
                    <MaskedTextField
                      data-testid="civilianTDYUBAllowance"
                      warning={civilianTDYUBAllowanceWarning}
                      defaultValue="0"
                      name="civilian_tdy_ub_allowance"
                      id="civilianTDYUBAllowance"
                      mask={Number}
                      scale={0}
                      signed={false}
                      thousandsSeparator=","
                      lazy={false}
                      labelHint={<span className={styles.civilianUBAllowanceWarning}>Optional</span>}
                      label={
                        <Label onClick={toggleCivilianTDYUBTooltip} className={styles.labelwithToolTip}>
                          If your orders specify a UB weight allowance, enter it here.
                          <ToolTip
                            text={
                              <span className={styles.toolTipText}>
                                If you do not specify a UB weight allowance, the default of 0 lbs will be used.
                              </span>
                            }
                            position="left"
                            icon="info-circle"
                            color="blue"
                            data-testid="civilianTDYUBAllowanceToolTip"
                            isVisible={showCivilianTDYUBTooltip}
                            closeOnLeave
                          />
                        </Label>
                      }
                    />
                  </div>
                </FormGroup>
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
    counseling_office_id: PropTypes.string,
    uploaded_orders: ExistingUploadsShape,
    civilian_tdy_ub_allowance: PropTypes.string,
  }).isRequired,
  onCancel: PropTypes.func.isRequired,
};

EditOrdersForm.defaultProps = {
  filePondEl: null,
};

const mapStateToProps = (state) => {
  return {
    affiliation: selectServiceMemberAffiliation(state) || '',
  };
};

const mapDispatchToProps = {
  setShowLoadingSpinner: setShowLoadingSpinnerAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(EditOrdersForm);
