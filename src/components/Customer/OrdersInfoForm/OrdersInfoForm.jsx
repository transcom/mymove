import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Radio, FormGroup, Label, Link as USWDSLink } from '@trussworks/react-uswds';
import { connect } from 'react-redux';

import { isBooleanFlagEnabled } from '../../../utils/featureFlags';
import { civilianTDYUBAllowanceWeightWarning, FEATURE_FLAG_KEYS } from '../../../shared/constants';

import styles from './OrdersInfoForm.module.scss';

import RequiredAsterisk, { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import ToolTip from 'shared/ToolTip/ToolTip';
import { ORDERS_PAY_GRADE_TYPE, ORDERS_TYPE } from 'constants/orders';
import { DropdownInput, DatePickerInput, DutyLocationInput } from 'components/form/fields';
import Hint from 'components/Hint/index';
import { Form } from 'components/form/Form';
import { DropdownArrayOf } from 'types';
import { DutyLocationShape } from 'types/dutyLocation';
import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import Callout from 'components/Callout';
import { formatLabelReportByDate, formatPayGradeOptions } from 'utils/formatters';
import { getPayGradeOptions, showCounselingOffices } from 'services/internalApi';
import { setShowLoadingSpinner as setShowLoadingSpinnerAction } from 'store/general/actions';
import retryPageLoading from 'utils/retryPageLoading';
import { milmoveLogger } from 'utils/milmoveLog';
import { selectServiceMemberAffiliation } from 'store/entities/selectors';

let originMeta;
let newDutyMeta = '';
const OrdersInfoForm = ({ ordersTypeOptions, initialValues, onSubmit, onBack, setShowLoadingSpinner, affiliation }) => {
  const [currentDutyLocation, setCurrentDutyLocation] = useState('');
  const [newDutyLocation, setNewDutyLocation] = useState('');
  const [counselingOfficeOptions, setCounselingOfficeOptions] = useState(null);
  const [showAccompaniedTourField, setShowAccompaniedTourField] = useState(false);
  const [showDependentAgeFields, setShowDependentAgeFields] = useState(false);
  const [hasDependents, setHasDependents] = useState(false);
  const [isOconusMove, setIsOconusMove] = useState(false);
  const [enableUB, setEnableUB] = useState(false);
  const [ordersType, setOrdersType] = useState('');
  const [grade, setGrade] = useState('');
  const [isCivilianTDYMove, setIsCivilianTDYMove] = useState(false);
  const [showCivilianTDYUBTooltip, setShowCivilianTDYUBTooltip] = useState(false);

  const [isHasDependentsDisabled, setHasDependentsDisabled] = useState(false);
  const [prevOrderType, setPrevOrderType] = useState('');
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
    grade: Yup.string().required('Required'),
    origin_duty_location: Yup.object().nullable().required('Required'),
    counseling_office_id: currentDutyLocation.provides_services_counseling
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

  useEffect(() => {
    // Functional component version of "componentDidMount"
    // By leaving the dependency array empty this will only run once
    const checkUBFeatureFlag = async () => {
      const enabled = await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.UNACCOMPANIED_BAGGAGE);
      if (enabled) {
        setEnableUB(true);
      }
    };
    checkUBFeatureFlag();
  }, []);

  useEffect(() => {
    const fetchCounselingOffices = async () => {
      if (currentDutyLocation?.id && !counselingOfficeOptions) {
        setShowLoadingSpinner(true, 'Loading counseling offices');
        try {
          const fetchedData = await showCounselingOffices(currentDutyLocation.id);
          if (fetchedData.body) {
            const counselingOffices = fetchedData.body.map((item) => ({
              key: item.id,
              value: item.name,
            }));
            setCounselingOfficeOptions(counselingOffices);
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
  }, [counselingOfficeOptions, currentDutyLocation.id, setShowLoadingSpinner]);

  const [payGradeOptions, setPayGradeOptions] = useState([]);
  useEffect(() => {
    const fetchGradeOptions = async () => {
      setShowLoadingSpinner(true, 'Loading Pay Grade options');
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
  }, [currentDutyLocation, newDutyLocation, isOconusMove, hasDependents, enableUB]);

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

  return (
    <Formik
      initialValues={initialValues}
      validateOnMount
      validationSchema={validationSchema}
      onSubmit={onSubmit}
      setShowAccompaniedTourField={setShowAccompaniedTourField}
      setShowDependentAgeFields={setShowDependentAgeFields}
    >
      {({ isValid, isSubmitting, handleSubmit, handleChange, setValues, values, touched, setFieldValue }) => {
        const isRetirementOrSeparation = ['RETIREMENT', 'SEPARATION'].includes(values.orders_type);

        const handleCounselingOfficeChange = () => {
          setValues({
            ...values,
            counseling_office_id: null,
          });
          setCounselingOfficeOptions(null);
        };
        if (!values.origin_duty_location && touched.origin_duty_location) originMeta = 'Required';
        else originMeta = null;

        if (!values.new_duty_location && touched.new_duty_location) newDutyMeta = 'Required';
        else newDutyMeta = null;

        const handleHasDependentsChange = (e) => {
          // Declare a duplicate local scope of the field value
          // for the form to prevent state race conditions
          if (e.target.value === '') {
            setFieldValue('has_dependents', '');
          } else {
            const fieldValueHasDependents = e.target.value === 'yes';
            setHasDependents(e.target.value === 'yes');
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
          <Form className={`${formStyles.form} ${styles.OrdersInfoForm}`}>
            <h1>Tell us about your move orders</h1>

            <SectionWrapper className={formStyles.formSection}>
              {requiredAsteriskMessage}
              <DropdownInput
                label="Orders type"
                name="orders_type"
                options={filteredOrderTypeOptions}
                required
                showRequiredAsterisk
                onChange={(e) => {
                  handleChange(e);
                  handleOrderTypeChange(e);
                }}
              />
              <DatePickerInput
                name="issue_date"
                label="Orders date"
                required
                showRequiredAsterisk
                renderInput={(input) => (
                  <>
                    {input}
                    <Hint>
                      <p>Date your orders were issued.</p>
                    </Hint>
                  </>
                )}
              />
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
                onDutyLocationChange={(e) => {
                  setCurrentDutyLocation(e);
                  handleCounselingOfficeChange();
                }}
                required
                showRequiredAsterisk
                metaOverride={originMeta}
              />
              {currentDutyLocation.provides_services_counseling && (
                <div>
                  <DropdownInput
                    label="Counseling office"
                    name="counseling_office_id"
                    id="counseling_office_id"
                    required
                    showRequiredAsterisk
                    options={counselingOfficeOptions}
                  />
                  <Hint>
                    Select an origin duty location that most closely represents your current physical location, not
                    where your shipment will originate, if different. This will allow a nearby transportation office to
                    assist you.
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
                    showRequiredAsterisk
                    displayAddress={false}
                    hint="Enter the option closest to your destination. Your move counselor will identify if there might be a cost to you."
                    metaOverride={newDutyMeta}
                    placeholder="Enter a city or ZIP"
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
                    showRequiredAsterisk
                    mask={Number}
                    scale={0}
                    signed={false}
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
                showRequiredAsterisk
                options={payGradeOptions}
                onChange={(e) => {
                  setGrade(e.target.value);
                  handleChange(e);
                }}
              />

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
                onBackClick={onBack}
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

OrdersInfoForm.propTypes = {
  ordersTypeOptions: DropdownArrayOf.isRequired,
  initialValues: PropTypes.shape({
    orders_type: PropTypes.string,
    issue_date: PropTypes.string,
    report_by_date: PropTypes.string,
    has_dependents: PropTypes.string,
    new_duty_location: DutyLocationShape,
    grade: PropTypes.string,
    origin_duty_location: DutyLocationShape,
    dependents_under_twelve: PropTypes.string,
    dependents_twelve_and_over: PropTypes.string,
    accompanied_tour: PropTypes.string,
    counseling_office_id: PropTypes.string,
    civilian_tdy_ub_allowance: PropTypes.string,
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
};

const mapDispatchToProps = {
  setShowLoadingSpinner: setShowLoadingSpinnerAction,
};

const mapStateToProps = (state) => {
  return {
    affiliation: selectServiceMemberAffiliation(state) || '',
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(OrdersInfoForm);
