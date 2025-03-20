import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Field, Formik } from 'formik';
import * as Yup from 'yup';
import { FormGroup, Label, Radio, Link as USWDSLink } from '@trussworks/react-uswds';
import { faInfoCircle } from '@fortawesome/free-solid-svg-icons';

import { isBooleanFlagEnabled } from '../../../utils/featureFlags';
import { FEATURE_FLAG_KEYS } from '../../../shared/constants';

import styles from './AddOrdersForm.module.scss';

import ToolTip from 'shared/ToolTip/ToolTip';
import { DatePickerInput, DropdownInput, DutyLocationInput } from 'components/form/fields';
import { Form } from 'components/form/Form';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { ORDERS_PAY_GRADE_OPTIONS, ORDERS_PAY_GRADE_TYPE, ORDERS_TYPE } from 'constants/orders';
import { dropdownInputOptions } from 'utils/formatters';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import Callout from 'components/Callout';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import formStyles from 'styles/form.module.scss';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import { showCounselingOffices } from 'services/ghcApi';
import Hint from 'components/Hint';

let originMeta;
let newDutyMeta = '';
const AddOrdersForm = ({
  onSubmit,
  ordersTypeOptions,
  initialValues,
  onBack,
  isSafetyMoveSelected,
  isBluebarkMoveSelected,
}) => {
  const payGradeOptions = dropdownInputOptions(ORDERS_PAY_GRADE_OPTIONS);
  const [counselingOfficeOptions, setCounselingOfficeOptions] = useState(null);
  const [currentDutyLocation, setCurrentDutyLocation] = useState('');
  const [newDutyLocation, setNewDutyLocation] = useState('');
  const [showAccompaniedTourField, setShowAccompaniedTourField] = useState(false);
  const [showDependentAgeFields, setShowDependentAgeFields] = useState(false);
  const [hasDependents, setHasDependents] = useState(false);
  const [isOconusMove, setIsOconusMove] = useState(false);
  const [enableUB, setEnableUB] = useState(false);
  const [isHasDependentsDisabled, setHasDependentsDisabled] = useState(false);
  const [prevOrderType, setPrevOrderType] = useState('');
  const [filteredOrderTypeOptions, setFilteredOrderTypeOptions] = useState(ordersTypeOptions);
  const [ordersType, setOrdersType] = useState('');
  const [grade, setGrade] = useState('');
  const [isCivilianTDYMove, setIsCivilianTDYMove] = useState(false);
  const { customerId: serviceMemberId } = useParams();

  const validationSchema = Yup.object().shape({
    ordersType: Yup.mixed()
      .oneOf(ordersTypeOptions.map((i) => i.key))
      .required('Required'),
    issueDate: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),
    reportByDate: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),
    hasDependents: Yup.mixed().oneOf(['yes', 'no']).required('Required'),
    originDutyLocation: Yup.object().nullable().required('Required'),
    counselingOfficeId: currentDutyLocation.provides_services_counseling
      ? Yup.string().required('Required')
      : Yup.string().notRequired(),
    newDutyLocation: Yup.object().nullable().required('Required'),
    grade: Yup.mixed().oneOf(Object.keys(ORDERS_PAY_GRADE_OPTIONS)).required('Required'),
    accompaniedTour: showAccompaniedTourField
      ? Yup.mixed().oneOf(['yes', 'no']).required('Required')
      : Yup.string().notRequired(),
    dependentsUnderTwelve: showDependentAgeFields
      ? Yup.number().min(0).required('Required')
      : Yup.number().notRequired(),
    dependentsTwelveAndOver: showDependentAgeFields
      ? Yup.number().min(0).required('Required')
      : Yup.number().notRequired(),
    civilianTdyUbAllowance: isCivilianTDYMove
      ? Yup.number()
          .transform((value) => (Number.isNaN(value) ? 0 : value))
          .min(0, 'UB weight allowance must be 0 or more')
          .max(2000, 'UB weight allowance cannot exceed 2000 lbs.')
      : Yup.number().notRequired(),
  });

  useEffect(() => {
    const checkUBFeatureFlag = async () => {
      const enabled = await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.UNACCOMPANIED_BAGGAGE);
      if (enabled) {
        setEnableUB(true);
      }
    };
    checkUBFeatureFlag();
  }, []);

  useEffect(() => {
    if (currentDutyLocation?.id && serviceMemberId) {
      showCounselingOffices(currentDutyLocation.id, serviceMemberId).then((fetchedData) => {
        if (fetchedData.body) {
          const counselingOffices = fetchedData.body.map((item) => ({
            key: item.id,
            value: item.name,
          }));
          setCounselingOfficeOptions(counselingOffices);
        }
      });
    }
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
  }, [currentDutyLocation, newDutyLocation, isOconusMove, hasDependents, enableUB, serviceMemberId]);

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
    <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ values, isValid, isSubmitting, handleSubmit, handleChange, touched, setFieldValue, setValues }) => {
        const isRetirementOrSeparation = ['RETIREMENT', 'SEPARATION'].includes(values.ordersType);
        if (!values.origin_duty_location && touched.origin_duty_location) originMeta = 'Required';
        else originMeta = null;

        const handleCounselingOfficeChange = () => {
          setValues({
            ...values,
            counselingOfficeId: null,
          });
          setCounselingOfficeOptions(null);
        };

        if (!values.newDutyLocation && touched.newDutyLocation) newDutyMeta = 'Required';
        else newDutyMeta = null;
        const handleHasDependentsChange = (e) => {
          // Declare a duplicate local scope of the field value
          // for the form to prevent state race conditions
          if (e.target.value === '') {
            setFieldValue('hasDependents', '');
          } else {
            const fieldValueHasDependents = e.target.value === 'yes';
            setHasDependents(e.target.value === 'yes');
            setFieldValue('hasDependents', fieldValueHasDependents ? 'yes' : 'no');
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
          values.civilianTdyUbAllowance > 350 && values.civilianTdyUbAllowance <= 2000;

        const civilianTDYUBAllowanceWeightWarning =
          '350 lbs. is the maximum UB weight allowance for a civilian TDY move unless stated otherwise on your orders.';

        let civilianTDYUBAllowanceWarning = '';
        if (showcivilianTDYUBAllowanceWarning) {
          civilianTDYUBAllowanceWarning = civilianTDYUBAllowanceWeightWarning;
        }
        return (
          <Form className={`${formStyles.form}`}>
            <ConnectedFlashMessage />
            <h1>Tell us about the orders</h1>

            <SectionWrapper className={formStyles.formSection}>
              <DropdownInput
                label="Orders type"
                name="ordersType"
                options={filteredOrderTypeOptions}
                required
                onChange={(e) => {
                  handleChange(e);
                  handleOrderTypeChange(e);
                }}
                isDisabled={isSafetyMoveSelected || isBluebarkMoveSelected}
                hint="Required"
              />
              <DatePickerInput name="issueDate" label="Orders date" required hint="Required" />
              <DatePickerInput name="reportByDate" label="Report by date" required hint="Required" />

              <DutyLocationInput
                label="Current duty location"
                name="originDutyLocation"
                id="originDutyLocation"
                onDutyLocationChange={(e) => {
                  setCurrentDutyLocation(e);
                  handleCounselingOfficeChange();
                }}
                metaOverride={originMeta}
                required
                hint="Required"
              />
              {currentDutyLocation.provides_services_counseling && (
                <div>
                  <DropdownInput
                    label="Counseling office"
                    name="counselingOfficeId"
                    id="counselingOfficeId"
                    data-testid="counselingOfficeSelect"
                    hint="Required"
                    required
                    options={counselingOfficeOptions}
                  />
                  <Hint>
                    Select an origin duty location that most closely represents the customers current physical location,
                    not where their shipment will originate, if different. This will allow a nearby transportation
                    office to assist them.
                  </Hint>
                </div>
              )}

              {isRetirementOrSeparation ? (
                <>
                  <h3>Where are they entitled to move?</h3>
                  <Callout>
                    <span>The government will pay for their move to:</span>
                    <ul>
                      <li>Home of record (HOR)</li>
                      <li>Place entered active duty (PLEAD)</li>
                    </ul>
                    <p>
                      It might pay for a move to their Home of selection (HOS), anywhere in CONUS. Check their orders.
                    </p>
                    <p>
                      Read more about where they are entitled to move when leaving the military on{' '}
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
                    name="newDutyLocation"
                    label="HOR, PLEAD or HOS"
                    displayAddress={false}
                    placeholder="Enter a city or ZIP"
                    metaOverride={newDutyMeta}
                    hint="Required"
                    onDutyLocationChange={(e) => {
                      setNewDutyLocation(e);
                    }}
                  />
                </>
              ) : (
                <DutyLocationInput
                  name="newDutyLocation"
                  label="New duty location"
                  required
                  hint="Required"
                  metaOverride={newDutyMeta}
                  onDutyLocationChange={(e) => {
                    setNewDutyLocation(e);
                  }}
                />
              )}

              <FormGroup>
                <Label hint="Required">Are dependents included in the orders?</Label>
                <div>
                  <Field
                    as={Radio}
                    label="Yes"
                    id="hasDependentsYes"
                    data-testid="hasDependentsYes"
                    name="hasDependents"
                    value="yes"
                    title="Yes, dependents are included in the orders"
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
                    name="hasDependents"
                    value="no"
                    title="No, dependents are not included in the orders"
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
                        name="accompaniedTour"
                        value="yes"
                        type="radio"
                      />
                      <ToolTip
                        text="Accompanied Tour: An authorized order (assignment or tour) that allows dependents to travel to the new Permanent Duty Station (PDS)"
                        position="right"
                        icon={faInfoCircle}
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
                        name="accompaniedTour"
                        value="no"
                        type="radio"
                      />
                      <ToolTip
                        text="Unaccompanied Tour: An authorized order (assignment or tour) that DOES NOT allow dependents to travel to the new Permanent Duty Station (PDS)"
                        position="right"
                        icon={faInfoCircle}
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
                    name="dependentsUnderTwelve"
                    label="Number of dependents under the age of 12"
                    id="dependentsUnderTwelve"
                    mask={Number}
                    scale={0}
                    signed={false}
                    thousandsSeparator=","
                    lazy={false}
                  />

                  <MaskedTextField
                    data-testid="dependentsTwelveAndOver"
                    defaultValue="0"
                    name="dependentsTwelveAndOver"
                    label="Number of dependents of the age 12 or over"
                    id="dependentsTwelveAndOver"
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
                options={payGradeOptions}
                hint="Required"
                onChange={(e) => {
                  setGrade(e.target.value);
                  handleChange(e);
                }}
              />

              {isCivilianTDYMove && (
                <FormGroup>
                  <div>
                    <MaskedTextField
                      data-testid="civilianTdyUbAllowance"
                      warning={
                        <span className={styles.civilianUBAllowanceWarning}>{civilianTDYUBAllowanceWarning}</span>
                      }
                      defaultValue="0"
                      name="civilianTdyUbAllowance"
                      id="civilianTdyUbAllowance"
                      mask={Number}
                      scale={0}
                      signed={false}
                      thousandsSeparator=","
                      lazy={false}
                      labelHint="Optional"
                      label={
                        <span className={styles.labelwithToolTip}>
                          If the customer&apos;s orders specify a specific UB weight allowance, enter it here.
                          <ToolTip
                            text="If you do not specify a UB weight allowance, the default of  0 lbs will be used."
                            position="right"
                            icon="info-circle"
                            color="blue"
                            data-testid="civilianTDYUBAllowanceToolTip"
                            closeOnLeave
                          />
                        </span>
                      }
                    />
                  </div>
                </FormGroup>
              )}
            </SectionWrapper>

            <div className={formStyles.formActions}>
              <WizardNavigation
                disableNext={!isValid || isSubmitting}
                onNextClick={handleSubmit}
                onBackClick={onBack}
              />
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

export default AddOrdersForm;
