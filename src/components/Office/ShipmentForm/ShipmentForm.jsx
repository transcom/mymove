import React, { useEffect, useState } from 'react';
import { arrayOf, bool, func, number, shape, string, oneOf } from 'prop-types';
import { Field, Formik } from 'formik';
import { generatePath, useNavigate, useParams, Link } from 'react-router-dom';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { Alert, Button, Checkbox, Fieldset, FormGroup, Radio, Label, Tag } from '@trussworks/react-uswds';
import classNames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import moment from 'moment';

import getShipmentOptions from '../../Customer/MtoShipmentForm/getShipmentOptions';
import { CloseoutOfficeInput } from '../../form/fields/CloseoutOfficeInput';

import ppmShipmentSchema from './ppmShipmentSchema';
import styles from './ShipmentForm.module.scss';
import MobileHomeShipmentForm from './MobileHomeShipmentForm/MobileHomeShipmentForm';
import mobileHomeShipmentSchema from './MobileHomeShipmentForm/mobileHomeShipmentSchema';
import BoatShipmentForm from './BoatShipmentForm/BoatShipmentForm';
import boatShipmentSchema from './BoatShipmentForm/boatShipmentSchema';

import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';
import SITCostDetails from 'components/Office/SITCostDetails/SITCostDetails';
import Hint from 'components/Hint/index';
import ConnectedDestructiveShipmentConfirmationModal from 'components/ConfirmationModals/DestructiveShipmentConfirmationModal';
import ConnectedShipmentAddressUpdateReviewRequestModal from 'components/Office/ShipmentAddressUpdateReviewRequestModal/ShipmentAddressUpdateReviewRequestModal';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { DatePickerInput, DropdownInput } from 'components/form/fields';
import { Form } from 'components/form';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentAccountingCodes from 'components/Office/ShipmentAccountingCodes/ShipmentAccountingCodes';
import ShipmentCustomerSIT from 'components/Office/ShipmentCustomerSIT/ShipmentCustomerSIT';
import ShipmentFormRemarks from 'components/Office/ShipmentFormRemarks/ShipmentFormRemarks';
import ShipmentIncentiveAdvance from 'components/Office/ShipmentIncentiveAdvance/ShipmentIncentiveAdvance';
import ShipmentVendor from 'components/Office/ShipmentVendor/ShipmentVendor';
import ShipmentWeight from 'components/Office/ShipmentWeight/ShipmentWeight';
import ShipmentWeightInput from 'components/Office/ShipmentWeightInput/ShipmentWeightInput';
import StorageFacilityAddress from 'components/Office/StorageFacilityAddress/StorageFacilityAddress';
import StorageFacilityInfo from 'components/Office/StorageFacilityInfo/StorageFacilityInfo';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { MOVES, MTO_SHIPMENTS } from 'constants/queryKeys';
import { servicesCounselingRoutes, tooRoutes } from 'constants/routes';
import { ADDRESS_UPDATE_STATUS, shipmentDestinationTypes } from 'constants/shipments';
import { officeRoles, roleTypes } from 'constants/userRoles';
import {
  deleteShipment,
  reviewShipmentAddressUpdate,
  updateMoveCloseoutOffice,
  dateSelectionIsWeekendHoliday,
} from 'services/ghcApi';
import {
  FEATURE_FLAG_KEYS,
  getPPMTypeLabel,
  PPM_TYPES,
  SHIPMENT_OPTIONS,
  SHIPMENT_TYPES,
  technicalHelpDeskURL,
} from 'shared/constants';
import formStyles from 'styles/form.module.scss';
import { AccountingCodesShape } from 'types/accountingCodes';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { ShipmentShape } from 'types/shipment';
import { TransportationOfficeShape } from 'types/transportationOffice';
import {
  formatMtoShipmentForAPI,
  formatMtoShipmentForDisplay,
  formatPpmShipmentForAPI,
  formatPpmShipmentForDisplay,
  formatMobileHomeShipmentForDisplay,
  formatMobileHomeShipmentForAPI,
  formatBoatShipmentForDisplay,
  formatBoatShipmentForAPI,
} from 'utils/formatMtoShipment';
import { formatWeight, dropdownInputOptions } from 'utils/formatters';
import { validateDate } from 'utils/validation';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { dateSelectionWeekendHolidayCheck } from 'utils/calendar';
import { datePickerFormat, formatDate, formatDateWithUTC } from 'shared/dates';
import { isPreceedingAddressComplete, isPreceedingAddressPPMPrimaryDestinationComplete } from 'shared/utils';
import { ORDERS_PAY_GRADE_TYPE } from 'constants/orders';
import { handleAddressToggleChange, blankAddress } from 'utils/shipments';
import { getResponseError } from 'services/internalApi';

const ShipmentForm = (props) => {
  const {
    newDutyLocationAddress,
    shipmentType,
    isCreatePage,
    isForServicesCounseling,
    mtoShipment,
    submitHandler,
    onUpdate,
    mtoShipments,
    serviceMember,
    currentResidence,
    moveTaskOrderID,
    TACs,
    SACs,
    userRole,
    displayDestinationType,
    isAdvancePage,
    move,
  } = props;

  const [estimatedWeightValue, setEstimatedWeightValue] = useState(mtoShipment?.ppmShipment?.estimatedWeight || 0);

  const updateEstimatedWeightValue = (value) => {
    setEstimatedWeightValue(value);
  };

  const { moveCode } = useParams();
  const navigate = useNavigate();

  const [datesErrorMessage, setDatesErrorMessage] = useState(null);
  const [errorMessage, setErrorMessage] = useState(null);
  const [errorCode, setErrorCode] = useState(null);
  const [successMessage, setSuccessMessage] = useState(null);
  const [shipmentAddressUpdateReviewErrorMessage, setShipmentAddressUpdateReviewErrorMessage] = useState(null);
  const [isCancelModalVisible, setIsCancelModalVisible] = useState(false);
  const [isAddressChangeModalOpen, setIsAddressChangeModalOpen] = useState(false);
  const [isTertiaryAddressEnabled, setIsTertiaryAddressEnabled] = useState(false);
  const [ppmSprFF, setPpmSprFF] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      setIsTertiaryAddressEnabled(await isBooleanFlagEnabled('third_address_available'));
    };
    fetchData();
  }, []);
  useEffect(() => {
    const fetchData = async () => {
      setPpmSprFF(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.PPM_SPR));
    };
    fetchData();
  }, []);

  const shipments = mtoShipments;

  const [isRequestedPickupDateAlertVisible, setIsRequestedPickupDateAlertVisible] = useState(false);
  const [isRequestedDeliveryDateAlertVisible, setIsRequestedDeliveryDateAlertVisible] = useState(false);
  const [requestedPickupDateAlertMessage, setRequestedPickupDateAlertMessage] = useState('');
  const [requestedDeliveryDateAlertMessage, setRequestedDeliveryDateAlertMessage] = useState('');
  const [isRequestedPickupDateInvalid, setIsRequestedPickupDateInvalid] = useState(false);
  const [isRequestedPickupDateChanged, setIsRequestedPickupDateChanged] = useState(false);
  const DEFAULT_COUNTRY_CODE = 'US';

  const queryClient = useQueryClient();
  const { mutate: mutateMTOShipmentStatus } = useMutation(deleteShipment, {
    onSuccess: (_, variables) => {
      const updatedMTOShipment = mtoShipment;
      // Update mtoShipments with our updated status and set query data to match
      shipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      // InvalidateQuery tells other components using this data that they need to re-fetch
      // This allows the requestCancellation button to update immediately
      queryClient.invalidateQueries([MTO_SHIPMENTS, variables.moveTaskOrderID]);

      // go back
      navigate(-1);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      setErrorMessage(errorMsg);
    },
  });

  const { mutate: mutateMoveCloseoutOffice } = useMutation(updateMoveCloseoutOffice, {
    onSuccess: () => {
      queryClient.invalidateQueries([MOVES, moveCode]);
    },
  });

  const { mutate: mutateShipmentAddressUpdateReview } = useMutation(reviewShipmentAddressUpdate, {
    onSuccess: (_, { successCallback }) => {
      setSuccessMessage('Changes sent to contractor.');
      setShipmentAddressUpdateReviewErrorMessage(null);
      setIsAddressChangeModalOpen(false);
      // After successfully updating, re-fetch MTO Shipments to get the shipment's updated address change request status
      queryClient
        .invalidateQueries([MTO_SHIPMENTS, moveTaskOrderID])
        .then(() => queryClient.refetchQueries([MTO_SHIPMENTS, moveTaskOrderID]));
      successCallback();
    },
    onError: () => {
      setSuccessMessage(null);
      setShipmentAddressUpdateReviewErrorMessage(
        'Something went wrong, and your changes were not saved. Please refresh the page and try again.',
      );
    },
  });

  const getShipmentNumber = () => {
    const { search } = window.location;
    const params = new URLSearchParams(search);
    const shipmentNumber = params.get('shipmentNumber');
    return shipmentNumber;
  };

  const handleDeleteShipment = (shipmentID) => {
    mutateMTOShipmentStatus({
      shipmentID,
    });
  };

  const handleSetError = (error, defaultErrorMessage) => {
    const { response } = error;

    if (setErrorCode && (response?.statusCode || response?.status)) {
      setErrorCode(response?.statusCode || response?.status);
    }

    const message = getResponseError(response, defaultErrorMessage);
    setErrorMessage(message);
  };

  const handleSubmitShipmentAddressUpdateReview = async (
    shipmentID,
    shipmentETag,
    status,
    officeRemarks,
    successCallback,
  ) => {
    mutateShipmentAddressUpdateReview({
      shipmentID,
      ifMatchETag: shipmentETag,
      body: {
        status,
        officeRemarks,
      },
      successCallback,
    });
  };

  const handleShowCancellationModal = () => {
    setIsCancelModalVisible(true);
  };

  // onload validate pickup date
  useEffect(() => {
    const onErrorHandler = (e) => {
      const { response } = e;
      setDatesErrorMessage(response?.body?.detail);
    };
    dateSelectionWeekendHolidayCheck(
      dateSelectionIsWeekendHoliday,
      DEFAULT_COUNTRY_CODE,
      new Date(mtoShipment.requestedPickupDate),
      'Requested pickup date',
      setRequestedPickupDateAlertMessage,
      setIsRequestedPickupDateAlertVisible,
      onErrorHandler,
    );
  }, [mtoShipment.requestedPickupDate]);

  // onload validate delivery date
  useEffect(() => {
    const onErrorHandler = (e) => {
      const { response } = e;
      setDatesErrorMessage(response?.body?.detail);
    };
    dateSelectionWeekendHolidayCheck(
      dateSelectionIsWeekendHoliday,
      DEFAULT_COUNTRY_CODE,
      new Date(mtoShipment.requestedDeliveryDate),
      'Requested delivery date',
      setRequestedDeliveryDateAlertMessage,
      setIsRequestedDeliveryDateAlertVisible,
      onErrorHandler,
    );
  }, [mtoShipment.requestedDeliveryDate]);

  const successMessageAlertControl = (
    <Button type="button" onClick={() => setSuccessMessage(null)} unstyled>
      <FontAwesomeIcon icon="times" className={styles.alertClose} />
    </Button>
  );

  const deliveryAddressUpdateRequested = mtoShipment?.deliveryAddressUpdate?.status === ADDRESS_UPDATE_STATUS.REQUESTED;

  const isHHG = shipmentType === SHIPMENT_OPTIONS.HHG;
  const isNTS = shipmentType === SHIPMENT_OPTIONS.NTS;
  const isNTSR = shipmentType === SHIPMENT_OPTIONS.NTSR;
  const isPPM = shipmentType === SHIPMENT_OPTIONS.PPM;
  const isMobileHome = shipmentType === SHIPMENT_OPTIONS.MOBILE_HOME;
  const isBoat =
    shipmentType === SHIPMENT_OPTIONS.BOAT ||
    shipmentType === SHIPMENT_TYPES.BOAT_HAUL_AWAY ||
    shipmentType === SHIPMENT_TYPES.BOAT_TOW_AWAY;
  const isUB = shipmentType === SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE;

  const showAccountingCodes = isNTS || isNTSR;

  const isTOO = userRole === roleTypes.TOO;
  const isServiceCounselor = userRole === roleTypes.SERVICES_COUNSELOR;
  const showCloseoutOffice =
    (isServiceCounselor || isTOO) &&
    isPPM &&
    (serviceMember.agency === SERVICE_MEMBER_AGENCIES.ARMY ||
      serviceMember.agency === SERVICE_MEMBER_AGENCIES.AIR_FORCE ||
      serviceMember.agency === SERVICE_MEMBER_AGENCIES.SPACE_FORCE);

  const shipmentDestinationAddressOptions = dropdownInputOptions(shipmentDestinationTypes);

  const shipmentNumber = isHHG ? getShipmentNumber() : null;
  let initialValues = {};
  if (isPPM) {
    initialValues = formatPpmShipmentForDisplay(
      isCreatePage
        ? { closeoutOffice: move.closeoutOffice }
        : {
            counselorRemarks: mtoShipment.counselorRemarks,
            ppmShipment: mtoShipment.ppmShipment,
            closeoutOffice: move.closeoutOffice,
          },
    );
  } else if (isMobileHome) {
    const hhgInitialValues = formatMtoShipmentForDisplay(
      isCreatePage ? { userRole } : { userRole, shipmentType, agents: mtoShipment.mtoAgents, ...mtoShipment },
    );
    initialValues = formatMobileHomeShipmentForDisplay(mtoShipment?.mobileHomeShipment, hhgInitialValues);
  } else if (isBoat) {
    const hhgInitialValues = formatMtoShipmentForDisplay(
      isCreatePage ? { userRole } : { userRole, shipmentType, agents: mtoShipment.mtoAgents, ...mtoShipment },
    );
    initialValues = formatBoatShipmentForDisplay(mtoShipment?.boatShipment, hhgInitialValues);
  } else {
    initialValues = formatMtoShipmentForDisplay(
      isCreatePage
        ? { userRole, shipmentType }
        : { userRole, shipmentType, agents: mtoShipment.mtoAgents, ...mtoShipment },
    );
  }

  let showDeliveryFields;
  let showPickupFields;
  let schema;

  if (isPPM) {
    schema = ppmShipmentSchema({
      estimatedIncentive: initialValues.estimatedIncentive || 0,
      weightAllotment: serviceMember.weightAllotment,
      advanceAmountRequested: mtoShipment.ppmShipment?.advanceAmountRequested,
      hasRequestedAdvance: mtoShipment.ppmShipment?.hasRequestedAdvance,
      isAdvancePage,
      showCloseoutOffice,
      sitEstimatedWeightMax: estimatedWeightValue || 0,
    });
  } else if (isMobileHome) {
    schema = mobileHomeShipmentSchema();
    showDeliveryFields = true;
    showPickupFields = true;
  } else if (isBoat) {
    schema = boatShipmentSchema();
    showDeliveryFields = true;
    showPickupFields = true;
  } else {
    const shipmentOptions = getShipmentOptions(shipmentType, userRole);

    showDeliveryFields = shipmentOptions.showDeliveryFields;
    showPickupFields = shipmentOptions.showPickupFields;
    schema = shipmentOptions.schema;
  }

  const optionalLabel = <span className={formStyles.optional}>Optional</span>;

  const moveDetailsPath = isTOO
    ? generatePath(tooRoutes.BASE_MOVE_VIEW_PATH, { moveCode })
    : generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, { moveCode });
  const editOrdersPath = isTOO
    ? generatePath(tooRoutes.BASE_ORDERS_EDIT_PATH, { moveCode })
    : generatePath(servicesCounselingRoutes.BASE_ORDERS_EDIT_PATH, { moveCode });

  const submitMTOShipment = (formValues, actions) => {
    //* PPM Shipment *//
    const newFormValues = formValues;
    if (isPPM) {
      if (isAdvancePage) {
        // Delete all form values not related to advance in order to avoid duplicate move history logs and unnecessary duplicate updates
        delete newFormValues.pickupAddress;
        delete newFormValues.destinationAddress;
        delete newFormValues.expectedDepartureDate;
      }
      const ppmShipmentBody = formatPpmShipmentForAPI(newFormValues);

      // Allow blank values to be entered into Pro Gear input fields
      if (
        ppmShipmentBody.ppmShipment.hasProGear &&
        ppmShipmentBody.ppmShipment.spouseProGearWeight >= 0 &&
        ppmShipmentBody.ppmShipment.proGearWeight === undefined
      ) {
        ppmShipmentBody.ppmShipment.proGearWeight = 0;
      }
      if (ppmShipmentBody.ppmShipment.hasProGear && ppmShipmentBody.ppmShipment.spouseProGearWeight === undefined) {
        ppmShipmentBody.ppmShipment.spouseProGearWeight = 0;
      }

      // Add a PPM shipment
      if (isCreatePage) {
        const body = { ...ppmShipmentBody, moveTaskOrderID };
        submitHandler(
          { body, normalize: false },
          {
            onSuccess: (newMTOShipment) => {
              const moveViewPath = generatePath(tooRoutes.BASE_MOVE_VIEW_PATH, { moveCode });
              const currentPath = isTOO
                ? generatePath(tooRoutes.BASE_SHIPMENT_EDIT_PATH, {
                    moveCode,
                    shipmentId: newMTOShipment.id,
                  })
                : generatePath(servicesCounselingRoutes.BASE_SHIPMENT_EDIT_PATH, {
                    moveCode,
                    shipmentId: newMTOShipment.id,
                  });
              const advancePath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_ADVANCE_PATH, {
                moveCode,
                shipmentId: newMTOShipment.id,
              });
              navigate(currentPath, { replace: true });
              if (isTOO) {
                navigate(moveViewPath);
              } else {
                navigate(advancePath);
              }
            },
            onError: (error) => {
              actions.setSubmitting(false);
              handleSetError(error, `Something went wrong, and your changes were not saved. Please try again.`);
            },
          },
        );
        return;
      }
      // Edit a PPM Shipment
      const updatePPMPayload = {
        moveTaskOrderID,
        shipmentID: mtoShipment.id,
        ifMatchETag: mtoShipment.eTag,
        normalize: false,
        body: ppmShipmentBody,
        locator: move.locator,
        moveETag: move.eTag,
      };

      const tooAdvancePath = generatePath(tooRoutes.BASE_SHIPMENT_ADVANCE_PATH_TOO, {
        moveCode,
        shipmentId: mtoShipment.id,
      });
      const advancePath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_ADVANCE_PATH, {
        moveCode,
        shipmentId: mtoShipment.id,
      });
      const SCMoveViewPath = generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, { moveCode });
      const tooMoveViewPath = generatePath(tooRoutes.BASE_MOVE_VIEW_PATH, { moveCode });

      submitHandler(updatePPMPayload, {
        onSuccess: () => {
          if (!isAdvancePage && formValues.closeoutOffice.id) {
            // If we are on the first page and a closeout office is a part of the form, we must be an SC editing a
            // PPM shipment, so we should update the closeout office and redirect to the advance page upon success.
            mutateMoveCloseoutOffice(
              {
                locator: moveCode,
                ifMatchETag: move.eTag,
                body: { closeoutOfficeId: formValues.closeoutOffice.id },
              },
              {
                onSuccess: () => {
                  actions.setSubmitting(false);
                  setErrorMessage(null);
                  navigate(advancePath);
                  onUpdate('success');
                },
                onError: (error) => {
                  actions.setSubmitting(false);
                  handleSetError(error, `Something went wrong, and your changes were not saved. Please try again.`);
                },
              },
            );
          } else if (!isAdvancePage && isServiceCounselor) {
            // If we are on the first page, and we are an SC with no closeout office present, we should redirect
            // to the advance page.
            actions.setSubmitting(false);
            navigate(advancePath);
            onUpdate('success');
          } else if (isServiceCounselor) {
            // If we are on the second page as an SC, we submit and redirect to the SC move view path.
            navigate(SCMoveViewPath);
            onUpdate('success');
          } else if (!isAdvancePage && isTOO) {
            actions.setSubmitting(false);
            navigate(tooMoveViewPath);
            onUpdate('success');
          } else {
            navigate(tooAdvancePath);
            onUpdate('success');
          }
        },
        onError: (error) => {
          actions.setSubmitting(false);
          handleSetError(error, `Something went wrong, and your changes were not saved. Please try again.`);
        },
      });
      return;
    }

    //* MTO Shipments *//

    const {
      pickup,
      delivery,
      customerRemarks,
      counselorRemarks,
      hasSecondaryDelivery,
      hasSecondaryPickup,
      secondaryPickup,
      secondaryDelivery,
      hasTertiaryDelivery,
      hasTertiaryPickup,
      tertiaryPickup,
      tertiaryDelivery,
      ntsRecordedWeight,
      tacType,
      sacType,
      serviceOrderNumber,
      storageFacility,
      usesExternalVendor,
      destinationType,
    } = formValues;

    let nullableTacType = tacType;
    let nullableSacType = sacType;
    if (showAccountingCodes && !isCreatePage) {
      nullableTacType = typeof tacType === 'undefined' ? '' : tacType;
      nullableSacType = typeof sacType === 'undefined' ? '' : sacType;
    }

    let pendingMtoShipment = formatMtoShipmentForAPI({
      shipmentType,
      moveCode,
      customerRemarks,
      counselorRemarks,
      pickup,
      delivery,
      ntsRecordedWeight,
      tacType: nullableTacType,
      sacType: nullableSacType,
      serviceOrderNumber,
      storageFacility,
      usesExternalVendor,
      destinationType,
      hasSecondaryPickup: hasSecondaryPickup === 'true',
      secondaryPickup: hasSecondaryPickup === 'true' ? secondaryPickup : {},
      hasSecondaryDelivery: hasSecondaryDelivery === 'true',
      secondaryDelivery: hasSecondaryDelivery === 'true' ? secondaryDelivery : {},
      hasTertiaryPickup: hasTertiaryPickup === 'true',
      tertiaryPickup: hasTertiaryPickup === 'true' ? tertiaryPickup : {},
      hasTertiaryDelivery: hasTertiaryDelivery === 'true',
      tertiaryDelivery: hasTertiaryDelivery === 'true' ? tertiaryDelivery : {},
    });

    // Mobile Home Shipment
    if (isMobileHome) {
      const mobileHomeShipmentBody = formatMobileHomeShipmentForAPI(formValues);
      pendingMtoShipment = {
        ...pendingMtoShipment,
        ...mobileHomeShipmentBody,
      };
    }

    // Boat Shipment
    if (isBoat) {
      const boatShipmentBody = formatBoatShipmentForAPI(formValues);
      pendingMtoShipment = {
        ...pendingMtoShipment,
        ...boatShipmentBody,
      };
    }

    const updateMTOShipmentPayload = {
      moveTaskOrderID,
      shipmentID: mtoShipment.id,
      ifMatchETag: mtoShipment.eTag,
      normalize: false,
      body: pendingMtoShipment,
    };

    // Add a MTO Shipment
    if (isCreatePage) {
      const body = { ...pendingMtoShipment, moveTaskOrderID };
      submitHandler(
        { body, normalize: false },
        {
          onSuccess: () => {
            navigate(moveDetailsPath);
          },
          onError: (error) => {
            handleSetError(error, `Something went wrong, and your changes were not saved. Please try again.`);
          },
        },
      );
    }
    // Edit MTO as Service Counselor
    else if (isForServicesCounseling) {
      // error handling handled in parent components
      submitHandler(updateMTOShipmentPayload, {
        onSuccess: () => {
          navigate(moveDetailsPath);
          onUpdate('success');
        },
        onError: (error) => {
          handleSetError(error, `Something went wrong, and your changes were not saved. Please try again.`);
        },
      });
    }
    // Edit a MTO Shipment as TOO
    else {
      submitHandler(updateMTOShipmentPayload, {
        onSuccess: () => {
          navigate(moveDetailsPath);
        },
        onError: (error) => {
          handleSetError(error, `Something went wrong, and your changes were not saved. Please try again.`);
        },
      });
    }
  };

  return (
    <Formik
      initialValues={initialValues}
      validateOnMount
      validateOnBlur
      validateOnChange
      validationSchema={schema}
      onSubmit={submitMTOShipment}
    >
      {({ values, isValid, isSubmitting, setValues, handleSubmit, setFieldError, validateForm, ...formikProps }) => {
        const {
          ppmType,
          hasSecondaryDestination,
          hasTertiaryDestination,
          hasDeliveryAddress,
          hasSecondaryPickup,
          hasSecondaryDelivery,
          hasTertiaryPickup,
          hasTertiaryDelivery,
        } = values;

        const isCivilian = serviceMember?.grade === ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE;
        if (!ppmType) {
          const type = isCivilian ? PPM_TYPES.ACTUAL_EXPENSE : PPM_TYPES.INCENTIVE_BASED;
          setValues({
            ...values,
            ppmType: type,
          });
        }
        const lengthHasError = !!(
          (formikProps.touched.lengthFeet && formikProps.errors.lengthFeet === 'Required') ||
          (formikProps.touched.lengthInches && formikProps.errors.lengthFeet === 'Required')
        );
        const widthHasError = !!(
          (formikProps.touched.widthFeet && formikProps.errors.widthFeet === 'Required') ||
          (formikProps.touched.widthInches && formikProps.errors.widthFeet === 'Required')
        );
        const heightHasError = !!(
          (formikProps.touched.heightFeet && formikProps.errors.heightFeet === 'Required') ||
          (formikProps.touched.heightInches && formikProps.errors.heightFeet === 'Required')
        );
        const dimensionError = !!(
          (formikProps.touched.lengthFeet && formikProps.errors.lengthFeet?.includes('Dimensions')) ||
          (formikProps.touched.lengthInches && formikProps.errors.lengthFeet?.includes('Dimensions'))
        );
        if (formikProps.touched.lengthInches && !formikProps.touched.lengthFeet) {
          formikProps.setFieldTouched('lengthFeet', true);
        }
        if (formikProps.touched.widthInches && !formikProps.touched.widthFeet) {
          formikProps.setFieldTouched('widthFeet', true);
        }
        if (formikProps.touched.heightInches && !formikProps.touched.heightFeet) {
          formikProps.setFieldTouched('heightFeet', true);
        }
        // manually turn off 'required' error when page loads if field is empty.
        if (values.year === null && !formikProps.touched.year && formikProps.errors.year === 'Required') {
          setFieldError('year', null);
        }

        const updateAddressTouched = (e, fieldName, address) => {
          handleAddressToggleChange(e, values, setValues, address);
          formikProps.setFieldTouched(`${fieldName}.usPostRegionCitiesID`, true);
        };

        const handleUseCurrentResidenceChange = (e) => {
          const { checked } = e.target;
          if (checked) {
            // use current residence
            setValues(
              {
                ...values,
                pickup: {
                  ...values.pickup,
                  address: currentResidence,
                },
              },
              { shouldValidate: true },
            );
          } else {
            // Revert address
            setValues(
              {
                ...values,
                pickup: {
                  ...values.pickup,
                  address: blankAddress.address,
                },
              },
              { shouldValidate: true },
            );
          }
        };

        const validatePickupDate = (e) => {
          let error = validateDate(e);
          // requestedPickupDate must be in the future for non-PPM shipments
          const pickupDate = moment(formatDateWithUTC(e)).startOf('day');
          const today = moment().startOf('day');

          if (!error && isRequestedPickupDateChanged && !isPPM && !pickupDate.isAfter(today)) {
            setIsRequestedPickupDateInvalid(true);
            error = 'Requested pickup date must be in the future.';
          } else {
            setIsRequestedPickupDateInvalid(false);
          }

          return error;
        };

        const handlePickupDateChange = (e) => {
          setValues({
            ...values,
            pickup: {
              ...values.pickup,
              requestedDate: formatDate(e, datePickerFormat),
            },
          });

          setIsRequestedPickupDateChanged(true);

          const onErrorHandler = (errResponse) => {
            const { response } = errResponse;
            setDatesErrorMessage(response?.body?.detail);
          };
          if (!validatePickupDate(e)) {
            dateSelectionWeekendHolidayCheck(
              dateSelectionIsWeekendHoliday,
              DEFAULT_COUNTRY_CODE,
              new Date(e),
              'Requested pickup date',
              setRequestedPickupDateAlertMessage,
              setIsRequestedPickupDateAlertVisible,
              onErrorHandler,
            );
          }
        };

        const handleDeliveryDateChange = (e) => {
          setValues({
            ...values,
            delivery: {
              ...values.delivery,
              requestedDate: formatDate(e, datePickerFormat),
            },
          });
          const onErrorHandler = (errResponse) => {
            const { response } = errResponse;
            setDatesErrorMessage(response?.body?.detail);
          };
          dateSelectionWeekendHolidayCheck(
            dateSelectionIsWeekendHoliday,
            DEFAULT_COUNTRY_CODE,
            new Date(e),
            'Requested delivery date',
            setRequestedDeliveryDateAlertMessage,
            setIsRequestedDeliveryDateAlertVisible,
            onErrorHandler,
          );
        };

        return (
          <>
            <ConnectedDestructiveShipmentConfirmationModal
              isOpen={isCancelModalVisible}
              shipmentID={mtoShipment.id}
              onClose={setIsCancelModalVisible}
              onSubmit={handleDeleteShipment}
            />
            <ConnectedShipmentAddressUpdateReviewRequestModal
              isOpen={isAddressChangeModalOpen}
              onClose={() => setIsAddressChangeModalOpen(false)}
              shipment={mtoShipment}
              onSubmit={async (shipmentID, shipmentETag, status, officeRemarks) => {
                const successCallback = () => {
                  if (status === ADDRESS_UPDATE_STATUS.APPROVED) {
                    setValues({
                      ...values,
                      hasDeliveryAddress: 'true',
                      delivery: {
                        ...values.delivery,
                        address: mtoShipment.deliveryAddressUpdate.newAddress,
                      },
                    });
                  }
                };
                await handleSubmitShipmentAddressUpdateReview(
                  shipmentID,
                  shipmentETag,
                  status,
                  officeRemarks,
                  successCallback,
                );
              }}
              errorMessage={shipmentAddressUpdateReviewErrorMessage}
              setErrorMessage={setShipmentAddressUpdateReviewErrorMessage}
            />
            <NotificationScrollToTop dependency={datesErrorMessage} />
            {datesErrorMessage && (
              <Alert data-testid="datesErrorMessage" type="error" headingLevel="h4" heading="An error occurred">
                {datesErrorMessage}
              </Alert>
            )}
            <NotificationScrollToTop dependency={errorMessage} />
            {errorMessage && (
              <Alert data-testid="errorMessage" type="error" headingLevel="h4" heading="An error occurred">
                {errorCode === 400 ? (
                  <p>
                    {errorMessage} Please try again later, or contact the&nbsp;
                    <Link to={technicalHelpDeskURL} target="_blank" rel="noreferrer">
                      Technical Help Desk
                    </Link>
                    .
                  </p>
                ) : (
                  <p>{errorMessage}</p>
                )}
              </Alert>
            )}
            <NotificationScrollToTop dependency={successMessage} />
            {successMessage && (
              <Alert type="success" cta={successMessageAlertControl}>
                {successMessage}
              </Alert>
            )}
            {isTOO && mtoShipment.usesExternalVendor && (
              <Alert headingLevel="h4" type="warning">
                The GHC prime contractor is not handling the shipment. Information will not be automatically shared with
                the movers handling it.
              </Alert>
            )}
            {deliveryAddressUpdateRequested && (
              <Alert type="error" className={styles.alert}>
                Request needs review. <a href="#delivery-location">See delivery address to proceed.</a>
              </Alert>
            )}

            <div className={styles.ShipmentForm}>
              <div className={styles.headerWrapper}>
                <div>
                  <ShipmentTag shipmentType={shipmentType} shipmentNumber={shipmentNumber} />
                  {ppmType === PPM_TYPES.SMALL_PACKAGE && (
                    <Tag className={styles.tagInfo} data-testid="ppmTypeTag">
                      {getPPMTypeLabel(ppmType)}
                    </Tag>
                  )}
                  {ppmType === PPM_TYPES.ACTUAL_EXPENSE && (
                    <Tag className={styles.tagInfo} data-testid="actualExpenseReimbursementTag">
                      {getPPMTypeLabel(ppmType)}
                    </Tag>
                  )}

                  <h1>{isCreatePage ? 'Add' : 'Edit'} shipment details</h1>
                </div>
                {!isCreatePage && mtoShipment?.status !== 'APPROVED' && (
                  <Button
                    type="button"
                    onClick={() => {
                      handleShowCancellationModal();
                    }}
                    unstyled
                  >
                    Delete shipment
                  </Button>
                )}
              </div>
              {isPPM && !isAdvancePage && isServiceCounselor && isCreatePage && (
                <SectionWrapper className={classNames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                  <Alert type="warning" aria-live="polite" headingLevel="h4" data-testid="scPPMCreateWarning">
                    Creating a PPM as a Service Counselor will automatically approve the PPM shipment and send it to the
                    customer. Please ensure all information is correct. After creation you will not be able to edit the
                    shipment.
                  </Alert>
                </SectionWrapper>
              )}
              <SectionWrapper className={styles.weightAllowance}>
                {isUB ? (
                  <p data-testid="ubWeightAllowance">
                    <strong>UB Weight allowance: </strong>
                    {formatWeight(serviceMember.weightAllotment.unaccompaniedBaggageAllowance)}
                  </p>
                ) : (
                  <p data-testid="weightAllowance">
                    <strong>Weight allowance: </strong>
                    {formatWeight(serviceMember.weightAllotment.authorizedWeight)}
                  </p>
                )}
              </SectionWrapper>

              <Form className={formStyles.form}>
                {isTOO && !isHHG && !isPPM && !isBoat && !isMobileHome && !isUB && <ShipmentVendor />}

                {isNTSR && <ShipmentWeightInput userRole={userRole} />}

                {isMobileHome && (
                  <MobileHomeShipmentForm
                    lengthHasError={lengthHasError}
                    widthHasError={widthHasError}
                    heightHasError={heightHasError}
                    values={values}
                    setFieldTouched={formikProps.setFieldTouched}
                    setFieldError={setFieldError}
                    validateForm={validateForm}
                    dimensionError={dimensionError}
                  />
                )}

                {isBoat && (
                  <BoatShipmentForm
                    lengthHasError={lengthHasError}
                    widthHasError={widthHasError}
                    heightHasError={heightHasError}
                    values={values}
                    setFieldTouched={formikProps.setFieldTouched}
                    setFieldError={setFieldError}
                    validateForm={validateForm}
                    dimensionError={dimensionError}
                  />
                )}

                {showPickupFields && (
                  <SectionWrapper className={formStyles.formSection}>
                    <h3 className={styles.SectionHeaderExtraSpacing}>Pickup details</h3>
                    <Fieldset data-testid="requestedPickupDateFieldSet">
                      {isRequestedPickupDateAlertVisible && !isRequestedPickupDateInvalid && (
                        <Alert
                          type="warning"
                          aria-live="polite"
                          headingLevel="h4"
                          data-testid="requestedPickupDateAlert"
                        >
                          {requestedPickupDateAlertMessage}
                        </Alert>
                      )}
                      <DatePickerInput
                        name="pickup.requestedDate"
                        label="Requested pickup date"
                        id="requestedPickupDate"
                        validate={validatePickupDate}
                        onChange={handlePickupDateChange}
                      />
                    </Fieldset>
                    {!isNTSR && (
                      <>
                        <AddressFields
                          name="pickup.address"
                          legend="Pickup Address"
                          formikProps={formikProps}
                          render={(fields) => (
                            <>
                              <p>What address are the movers picking up from?</p>
                              <Checkbox
                                data-testid="useCurrentResidence"
                                label="Use pickup address"
                                name="useCurrentResidence"
                                onChange={handleUseCurrentResidenceChange}
                                id="useCurrentResidenceCheckbox"
                              />
                              {fields}
                              <h4>Second Pickup Address</h4>
                              <FormGroup>
                                <p>Do you want movers to pick up any belongings from a second address?</p>
                                <div className={formStyles.radioGroup}>
                                  <Field
                                    as={Radio}
                                    id="has-secondary-pickup"
                                    data-testid="has-secondary-pickup"
                                    label="Yes"
                                    name="hasSecondaryPickup"
                                    value="true"
                                    title="Yes, I have a second pickup address"
                                    checked={hasSecondaryPickup === 'true'}
                                    disabled={!isPreceedingAddressComplete('true', values.pickup.address)}
                                    onChange={(e) => updateAddressTouched(e, 'secondaryPickup.address', blankAddress)}
                                  />
                                  <Field
                                    as={Radio}
                                    id="no-secondary-pickup"
                                    data-testid="no-secondary-pickup"
                                    label="No"
                                    name="hasSecondaryPickup"
                                    value="false"
                                    title="No, I do not have a second pickup address"
                                    checked={hasSecondaryPickup !== 'true'}
                                    disabled={!isPreceedingAddressComplete('true', values.pickup.address)}
                                    onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                  />
                                </div>
                              </FormGroup>
                              {hasSecondaryPickup === 'true' && (
                                <>
                                  <AddressFields name="secondaryPickup.address" formikProps={formikProps} />
                                  {isTertiaryAddressEnabled && (
                                    <>
                                      <h4>Third Pickup Address</h4>
                                      <FormGroup>
                                        <p>Do you want movers to pick up any belongings from a third address?</p>
                                        <div className={formStyles.radioGroup}>
                                          <Field
                                            as={Radio}
                                            id="has-tertiary-pickup"
                                            data-testid="has-tertiary-pickup"
                                            label="Yes"
                                            name="hasTertiaryPickup"
                                            value="true"
                                            title="Yes, I have a third pickup address"
                                            checked={hasTertiaryPickup === 'true'}
                                            disabled={
                                              !isPreceedingAddressComplete(
                                                hasSecondaryPickup,
                                                values.secondaryPickup.address,
                                              )
                                            }
                                            onChange={(e) =>
                                              updateAddressTouched(e, 'tertiaryPickup.address', blankAddress)
                                            }
                                          />
                                          <Field
                                            as={Radio}
                                            id="no-tertiary-pickup"
                                            data-testid="no-tertiary-pickup"
                                            label="No"
                                            name="hasTertiaryPickup"
                                            value="false"
                                            title="No, I do not have a third pickup address"
                                            checked={hasTertiaryPickup !== 'true'}
                                            disabled={
                                              !isPreceedingAddressComplete(
                                                hasSecondaryPickup,
                                                values.secondaryPickup.address,
                                              )
                                            }
                                            onChange={(e) =>
                                              handleAddressToggleChange(e, values, setValues, blankAddress)
                                            }
                                          />
                                        </div>
                                      </FormGroup>
                                      {hasTertiaryPickup === 'true' && (
                                        <AddressFields name="tertiaryPickup.address" formikProps={formikProps} />
                                      )}
                                    </>
                                  )}
                                </>
                              )}
                            </>
                          )}
                        />

                        <ContactInfoFields
                          name="pickup.agent"
                          legend={<div className={formStyles.legendContent}>Releasing agent {optionalLabel}</div>}
                          render={(fields) => {
                            return fields;
                          }}
                        />
                      </>
                    )}
                  </SectionWrapper>
                )}

                {isTOO && (isNTS || isNTSR) && (
                  <>
                    <StorageFacilityInfo userRole={userRole} />
                    <StorageFacilityAddress values={values} formikProps={formikProps} />
                  </>
                )}

                {isServiceCounselor && isNTSR && (
                  <>
                    <StorageFacilityInfo userRole={userRole} />
                    <StorageFacilityAddress values={values} formikProps={formikProps} />
                  </>
                )}

                {showDeliveryFields && (
                  <SectionWrapper className={formStyles.formSection}>
                    <h3 className={styles.SectionHeaderExtraSpacing}>Delivery details</h3>
                    <Fieldset>
                      {isRequestedDeliveryDateAlertVisible && (
                        <Alert type="warning" aria-live="polite" headingLevel="h4">
                          {requestedDeliveryDateAlertMessage}
                        </Alert>
                      )}
                      <DatePickerInput
                        name="delivery.requestedDate"
                        label="Requested delivery date"
                        id="requestedDeliveryDate"
                        validate={validateDate}
                        onChange={handleDeliveryDateChange}
                      />
                    </Fieldset>
                    {isNTSR && (
                      <>
                        {deliveryAddressUpdateRequested && (
                          <Alert type="error" slim className={styles.deliveryAddressUpdateAlert} id="delivery-location">
                            <span className={styles.deliveryAddressUpdateAlertContent}>
                              Pending delivery address change request needs review.{' '}
                              <Button
                                className={styles.reviewRequestLink}
                                type="button"
                                unstyled
                                onClick={() => setIsAddressChangeModalOpen(true)}
                                disabled={false}
                              >
                                Review request
                              </Button>{' '}
                              to proceed.
                            </span>
                          </Alert>
                        )}
                        <Fieldset
                          legend="Delivery Address"
                          disabled={deliveryAddressUpdateRequested}
                          className={classNames('usa-legend')}
                        >
                          <AddressFields
                            name="delivery.address"
                            formikProps={formikProps}
                            render={(fields) => {
                              return fields;
                            }}
                          />
                          <h4>Second Delivery Address</h4>
                          <FormGroup>
                            <p>Do you want the movers to deliver any belongings to a second address?</p>
                            <div className={formStyles.radioGroup}>
                              <Field
                                as={Radio}
                                data-testid="has-secondary-delivery"
                                id="has-secondary-delivery"
                                label="Yes"
                                name="hasSecondaryDelivery"
                                value="true"
                                title="Yes, I have a second destination location"
                                checked={hasSecondaryDelivery === 'true'}
                                disabled={!isPreceedingAddressComplete('true', values.delivery.address)}
                                onChange={(e) => updateAddressTouched(e, 'secondaryDelivery.address', blankAddress)}
                              />
                              <Field
                                as={Radio}
                                data-testid="no-secondary-delivery"
                                id="no-secondary-delivery"
                                label="No"
                                name="hasSecondaryDelivery"
                                value="false"
                                title="No, I do not have a second destination location"
                                checked={hasSecondaryDelivery !== 'true'}
                                disabled={!isPreceedingAddressComplete('true', values.delivery.address)}
                                onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                              />
                            </div>
                          </FormGroup>
                          {hasSecondaryDelivery === 'true' && (
                            <>
                              <AddressFields name="secondaryDelivery.address" formikProps={formikProps} />
                              {isTertiaryAddressEnabled && (
                                <>
                                  <h4>Third Delivery Address</h4>
                                  <FormGroup>
                                    <p>Do you want the movers to deliver any belongings to a third address?</p>
                                    <div className={formStyles.radioGroup}>
                                      <Field
                                        as={Radio}
                                        id="has-tertiary-delivery"
                                        data-testid="has-tertiary-delivery"
                                        label="Yes"
                                        name="hasTertiaryDelivery"
                                        value="true"
                                        title="Yes, I have a third delivery address"
                                        checked={hasTertiaryDelivery === 'true'}
                                        disabled={
                                          !isPreceedingAddressComplete(
                                            hasSecondaryDelivery,
                                            values.secondaryDelivery.address,
                                          )
                                        }
                                        onChange={(e) =>
                                          updateAddressTouched(e, 'tertiaryDelivery.address', blankAddress)
                                        }
                                      />
                                      <Field
                                        as={Radio}
                                        id="no-tertiary-delivery"
                                        data-testid="no-tertiary-delivery"
                                        label="No"
                                        name="hasTertiaryDelivery"
                                        value="false"
                                        title="No, I do not have a third delivery address"
                                        checked={hasTertiaryDelivery !== 'true'}
                                        disabled={
                                          !isPreceedingAddressComplete(
                                            hasSecondaryDelivery,
                                            values.secondaryDelivery.address,
                                          )
                                        }
                                        onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                      />
                                    </div>
                                  </FormGroup>
                                  {hasTertiaryDelivery === 'true' && (
                                    <AddressFields name="tertiaryDelivery.address" formikProps={formikProps} />
                                  )}
                                </>
                              )}
                            </>
                          )}
                          {displayDestinationType && (
                            <DropdownInput
                              label="Destination type"
                              name="destinationType"
                              options={shipmentDestinationAddressOptions}
                              id="destinationType"
                            />
                          )}
                        </Fieldset>

                        <ContactInfoFields
                          name="delivery.agent"
                          legend={<div className={formStyles.legendContent}>Receiving agent {optionalLabel}</div>}
                          render={(fields) => {
                            return fields;
                          }}
                        />
                      </>
                    )}
                    {!isNTS && !isNTSR && (
                      <>
                        <p className={classNames('usa-legend', styles.mockLegend)} id="delivery-location">
                          Delivery Address
                        </p>
                        {deliveryAddressUpdateRequested && (
                          <Alert type="error" slim className={styles.deliveryAddressUpdateAlert}>
                            <span className={styles.deliveryAddressUpdateAlertContent}>
                              Pending delivery address change request needs review.{' '}
                              <Button
                                className={styles.reviewRequestLink}
                                type="button"
                                unstyled
                                onClick={() => setIsAddressChangeModalOpen(true)}
                                disabled={false}
                              >
                                Review request
                              </Button>{' '}
                              to proceed.
                            </span>
                          </Alert>
                        )}
                        <Fieldset
                          legendStyle="srOnly"
                          legend="Delivery Address"
                          disabled={deliveryAddressUpdateRequested}
                        >
                          <FormGroup>
                            <p>Does the customer know their delivery address yet?</p>
                            <div className={formStyles.radioGroup}>
                              <Field
                                as={Radio}
                                id="has-delivery-address"
                                label="Yes"
                                name="hasDeliveryAddress"
                                value="true"
                                title="Yes, I know my delivery address"
                                checked={hasDeliveryAddress === 'true'}
                                onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                              />
                              <Field
                                as={Radio}
                                id="no-delivery-address"
                                label="No"
                                name="hasDeliveryAddress"
                                value="false"
                                title="No, I do not know my delivery address"
                                checked={hasDeliveryAddress === 'false'}
                                onChange={(e) => updateAddressTouched(e, 'delivery.address', newDutyLocationAddress)}
                              />
                            </div>
                          </FormGroup>
                          {hasDeliveryAddress === 'true' ? (
                            <AddressFields
                              name="delivery.address"
                              formikProps={formikProps}
                              render={(fields) => (
                                <>
                                  {fields}
                                  {displayDestinationType && (
                                    <DropdownInput
                                      label="Destination type"
                                      name="destinationType"
                                      options={shipmentDestinationAddressOptions}
                                      id="destinationType"
                                    />
                                  )}
                                  <h4>Second Delivery Address</h4>
                                  <FormGroup>
                                    <p>Do you want the movers to deliver any belongings to a second address?</p>
                                    <div className={formStyles.radioGroup}>
                                      <Field
                                        as={Radio}
                                        data-testid="has-secondary-delivery"
                                        id="has-secondary-delivery"
                                        label="Yes"
                                        name="hasSecondaryDelivery"
                                        value="true"
                                        title="Yes, I have a second destination location"
                                        checked={hasSecondaryDelivery === 'true'}
                                        disabled={
                                          !isPreceedingAddressComplete(hasDeliveryAddress, values.delivery.address)
                                        }
                                        onChange={(e) =>
                                          updateAddressTouched(e, 'secondaryDelivery.address', blankAddress)
                                        }
                                      />
                                      <Field
                                        as={Radio}
                                        data-testid="no-secondary-delivery"
                                        id="no-secondary-delivery"
                                        label="No"
                                        name="hasSecondaryDelivery"
                                        value="false"
                                        title="No, I do not have a second destination location"
                                        checked={hasSecondaryDelivery !== 'true'}
                                        disabled={
                                          !isPreceedingAddressComplete(hasDeliveryAddress, values.delivery.address)
                                        }
                                        onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                      />
                                    </div>
                                  </FormGroup>
                                  {hasSecondaryDelivery === 'true' && (
                                    <>
                                      <AddressFields name="secondaryDelivery.address" formikProps={formikProps} />
                                      {isTertiaryAddressEnabled && (
                                        <>
                                          <h4>Third Delivery Address</h4>
                                          <FormGroup>
                                            <p>Do you want the movers to deliver any belongings to a third address?</p>
                                            <div className={formStyles.radioGroup}>
                                              <Field
                                                as={Radio}
                                                id="has-tertiary-delivery"
                                                data-testid="has-tertiary-delivery"
                                                label="Yes"
                                                name="hasTertiaryDelivery"
                                                value="true"
                                                title="Yes, I have a third delivery address"
                                                checked={hasTertiaryDelivery === 'true'}
                                                disabled={
                                                  !isPreceedingAddressComplete(
                                                    hasSecondaryDelivery,
                                                    values.secondaryDelivery.address,
                                                  )
                                                }
                                                onChange={(e) =>
                                                  updateAddressTouched(e, 'tertiaryDelivery.address', blankAddress)
                                                }
                                              />
                                              <Field
                                                as={Radio}
                                                id="no-tertiary-delivery"
                                                data-testid="no-tertiary-delivery"
                                                label="No"
                                                name="hasTertiaryDelivery"
                                                value="false"
                                                title="No, I do not have a third delivery address"
                                                checked={hasTertiaryDelivery !== 'true'}
                                                disabled={
                                                  !isPreceedingAddressComplete(
                                                    hasSecondaryDelivery,
                                                    values.secondaryDelivery.address,
                                                  )
                                                }
                                                onChange={(e) =>
                                                  handleAddressToggleChange(e, values, setValues, blankAddress)
                                                }
                                              />
                                            </div>
                                          </FormGroup>
                                          {hasTertiaryDelivery === 'true' && (
                                            <AddressFields name="tertiaryDelivery.address" formikProps={formikProps} />
                                          )}
                                        </>
                                      )}
                                    </>
                                  )}
                                </>
                              )}
                            />
                          ) : (
                            <div>
                              <p>
                                We can use the zip of their{' '}
                                {displayDestinationType ? 'HOR, HOS or PLEAD:' : 'new duty location:'}
                                <br />
                                <strong>
                                  {newDutyLocationAddress.city}, {newDutyLocationAddress.state}{' '}
                                  {newDutyLocationAddress.postalCode}{' '}
                                </strong>
                              </p>
                              {displayDestinationType && (
                                <DropdownInput
                                  label="Destination type"
                                  name="destinationType"
                                  options={shipmentDestinationAddressOptions}
                                  id="destinationType"
                                />
                              )}
                            </div>
                          )}
                        </Fieldset>

                        <ContactInfoFields
                          name="delivery.agent"
                          legend={<div className={formStyles.legendContent}>Receiving agent {optionalLabel}</div>}
                          render={(fields) => {
                            return fields;
                          }}
                        />
                      </>
                    )}
                  </SectionWrapper>
                )}

                {isPPM && !isAdvancePage && (
                  <>
                    {isServiceCounselor && (
                      <SectionWrapper
                        className={classNames(ppmStyles.sectionWrapper, formStyles.formSection)}
                        data-testid="ppmTypeSection"
                      >
                        <h3>PPM Type</h3>
                        <FormGroup>
                          <Label className={styles.Label} htmlFor="ppmType">
                            Indicate the PPM Type
                          </Label>
                          <Field
                            as={Radio}
                            id="isIncentiveBased"
                            label={getPPMTypeLabel(PPM_TYPES.INCENTIVE_BASED)}
                            name="ppmType"
                            value={PPM_TYPES.INCENTIVE_BASED}
                            checked={(ppmType == null && !isCivilian) || ppmType === PPM_TYPES.INCENTIVE_BASED}
                            disabled={isCivilian}
                            className={styles.buttonGroup}
                            data-testid="isIncentiveBased"
                          />
                          <Field
                            as={Radio}
                            id="isActualExpense"
                            label={getPPMTypeLabel(PPM_TYPES.ACTUAL_EXPENSE)}
                            name="ppmType"
                            value={PPM_TYPES.ACTUAL_EXPENSE}
                            checked={(ppmType == null && isCivilian) || ppmType === PPM_TYPES.ACTUAL_EXPENSE}
                            className={styles.buttonGroup}
                            data-testid="isActualExpense"
                          />
                          {ppmSprFF && (
                            <Field
                              as={Radio}
                              id="isSmallPackage"
                              label={getPPMTypeLabel(PPM_TYPES.SMALL_PACKAGE)}
                              name="ppmType"
                              value={PPM_TYPES.SMALL_PACKAGE}
                              checked={ppmType === PPM_TYPES.SMALL_PACKAGE}
                              className={styles.buttonGroup}
                              data-testid="isSmallPackage"
                            />
                          )}
                        </FormGroup>
                      </SectionWrapper>
                    )}
                    <SectionWrapper className={classNames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                      <h3>{ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Shipped Date' : 'Departure Date'}</h3>
                      <DatePickerInput
                        name="expectedDepartureDate"
                        label={
                          ppmType === PPM_TYPES.SMALL_PACKAGE
                            ? 'When did the customer ship their package?'
                            : 'Planned Departure Date'
                        }
                      />
                      <Hint className={ppmStyles.hint}>
                        Enter the first day you expect to move things. It&apos;s OK if the actual date is different. We
                        will ask for your actual departure date when you document and complete your PPM.
                      </Hint>
                    </SectionWrapper>
                    <SectionWrapper className={classNames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                      <AddressFields
                        name="pickup.address"
                        legend={ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Shipped from Address' : 'Pickup Address'}
                        formikProps={formikProps}
                        render={(fields) => (
                          <>
                            <p>What address are you moving from?</p>
                            <Checkbox
                              data-testid="useCurrentResidence"
                              label="Use Pickup Address"
                              name="useCurrentResidence"
                              onChange={handleUseCurrentResidenceChange}
                              id="useCurrentResidenceCheckbox"
                            />
                            {fields}
                            <h4>Second Pickup Address</h4>
                            <FormGroup>
                              <p>
                                Will you move any belongings from a second address? (Must be near the pickup address.
                                Subject to approval.)
                              </p>
                              <div className={formStyles.radioGroup}>
                                <Field
                                  as={Radio}
                                  id="has-secondary-pickup"
                                  data-testid="has-secondary-pickup"
                                  label="Yes"
                                  name="hasSecondaryPickup"
                                  value="true"
                                  title="Yes, there is a second pickup address"
                                  checked={hasSecondaryPickup === 'true'}
                                  disabled={!isPreceedingAddressComplete('true', values.pickup.address)}
                                  onChange={(e) => updateAddressTouched(e, 'secondaryPickup.address', blankAddress)}
                                />
                                <Field
                                  as={Radio}
                                  id="no-secondary-pickup"
                                  data-testid="no-secondary-pickup"
                                  label="No"
                                  name="hasSecondaryPickup"
                                  value="false"
                                  title="No, there is not a second pickup address"
                                  checked={hasSecondaryPickup !== 'true'}
                                  disabled={!isPreceedingAddressComplete('true', values.pickup.address)}
                                  onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                />
                              </div>
                            </FormGroup>
                            {hasSecondaryPickup === 'true' && (
                              <>
                                <AddressFields name="secondaryPickup.address" formikProps={formikProps} />
                                {isTertiaryAddressEnabled && (
                                  <>
                                    <h4>Third Pickup Address</h4>
                                    <FormGroup>
                                      <p>
                                        Will you move any belongings from a third address? (Must be near the pickup
                                        address. Subject to approval.)
                                      </p>
                                      <div className={formStyles.radioGroup}>
                                        <Field
                                          as={Radio}
                                          id="has-tertiary-pickup"
                                          data-testid="has-tertiary-pickup"
                                          label="Yes"
                                          name="hasTertiaryPickup"
                                          value="true"
                                          title="Yes, there is a third pickup address"
                                          checked={hasTertiaryPickup === 'true'}
                                          disabled={
                                            !isPreceedingAddressComplete(
                                              hasSecondaryPickup,
                                              values.secondaryPickup.address,
                                            )
                                          }
                                          onChange={(e) =>
                                            updateAddressTouched(e, 'tertiaryPickup.address', blankAddress)
                                          }
                                        />
                                        <Field
                                          as={Radio}
                                          id="no-tertiary-pickup"
                                          data-testid="no-tertiary-pickup"
                                          label="No"
                                          name="hasTertiaryPickup"
                                          value="false"
                                          title="No, there is not a third pickup address"
                                          checked={hasTertiaryPickup !== 'true'}
                                          disabled={
                                            !isPreceedingAddressComplete(
                                              hasSecondaryPickup,
                                              values.secondaryPickup.address,
                                            )
                                          }
                                          onChange={(e) =>
                                            handleAddressToggleChange(e, values, setValues, blankAddress)
                                          }
                                        />
                                      </div>
                                    </FormGroup>
                                    {hasTertiaryPickup === 'true' && (
                                      <AddressFields name="tertiaryPickup.address" formikProps={formikProps} />
                                    )}
                                  </>
                                )}
                              </>
                            )}
                          </>
                        )}
                      />
                      <AddressFields
                        name="destination.address"
                        legend={ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Destination Address' : 'Delivery Address'}
                        formikProps={formikProps}
                        address1LabelHint="Optional"
                        render={(fields) => (
                          <>
                            {fields}
                            <h4>Second {ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Destination' : 'Delivery'} Address</h4>
                            <FormGroup>
                              <p>
                                Will you move any belongings to a second address? (Must be near the{' '}
                                {ppmType === PPM_TYPES.SMALL_PACKAGE ? 'destination' : 'delivery'} address. Subject to
                                approval.)
                              </p>
                              <div className={formStyles.radioGroup}>
                                <Field
                                  as={Radio}
                                  data-testid="has-secondary-destination"
                                  id="has-secondary-destination"
                                  label="Yes"
                                  name="hasSecondaryDestination"
                                  value="true"
                                  title="Yes, there is a second destination location"
                                  checked={hasSecondaryDestination === 'true'}
                                  disabled={
                                    !isPreceedingAddressPPMPrimaryDestinationComplete(values.destination.address)
                                  }
                                  onChange={(e) =>
                                    updateAddressTouched(e, 'secondaryDestination.address', blankAddress)
                                  }
                                />
                                <Field
                                  as={Radio}
                                  data-testid="no-secondary-destination"
                                  id="no-secondary-destination"
                                  label="No"
                                  name="hasSecondaryDestination"
                                  value="false"
                                  title="No, there is not a second destination location"
                                  checked={hasSecondaryDestination !== 'true'}
                                  disabled={
                                    !isPreceedingAddressPPMPrimaryDestinationComplete(values.destination.address)
                                  }
                                  onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                />
                              </div>
                            </FormGroup>
                            {hasSecondaryDestination === 'true' && (
                              <>
                                <AddressFields name="secondaryDestination.address" formikProps={formikProps} />
                                {isTertiaryAddressEnabled && (
                                  <>
                                    <h4>
                                      Third {ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Destination' : 'Delivery'} Address
                                    </h4>
                                    <FormGroup>
                                      <p>
                                        Will you move any belongings to a third address? (Must be near the{' '}
                                        {ppmType === PPM_TYPES.SMALL_PACKAGE ? 'destination' : 'delivery'} address.
                                        Subject to approval.)
                                      </p>
                                      <div className={formStyles.radioGroup}>
                                        <Field
                                          as={Radio}
                                          id="has-tertiary-destination"
                                          data-testid="has-tertiary-destination"
                                          label="Yes"
                                          name="hasTertiaryDestination"
                                          value="true"
                                          title="Yes, I have a third delivery address"
                                          checked={hasTertiaryDestination === 'true'}
                                          disabled={
                                            !isPreceedingAddressComplete(
                                              hasSecondaryDestination,
                                              values.secondaryDestination.address,
                                            )
                                          }
                                          onChange={(e) =>
                                            updateAddressTouched(e, 'tertiaryDestination.address', blankAddress)
                                          }
                                        />
                                        <Field
                                          as={Radio}
                                          id="no-tertiary-destination"
                                          data-testid="no-tertiary-destination"
                                          label="No"
                                          name="hasTertiaryDestination"
                                          value="false"
                                          title="No, I do not have a third delivery address"
                                          checked={hasTertiaryDestination !== 'true'}
                                          disabled={
                                            !isPreceedingAddressComplete(
                                              hasSecondaryDestination,
                                              values.secondaryDestination.address,
                                            )
                                          }
                                          onChange={(e) =>
                                            handleAddressToggleChange(e, values, setValues, blankAddress)
                                          }
                                        />
                                      </div>
                                    </FormGroup>
                                    {hasTertiaryDestination === 'true' && (
                                      <AddressFields name="tertiaryDestination.address" formikProps={formikProps} />
                                    )}
                                  </>
                                )}
                              </>
                            )}
                          </>
                        )}
                      />
                    </SectionWrapper>
                    {showCloseoutOffice && (
                      <SectionWrapper>
                        <h3>Closeout office</h3>
                        <CloseoutOfficeInput
                          hint="If there is more than one PPM for this move, the closeout office will be the same for all your PPMs."
                          name="closeoutOffice"
                          placeholder="Start typing a closeout location..."
                          label="Closeout location"
                          displayAddress
                        />
                      </SectionWrapper>
                    )}
                    {ppmType !== PPM_TYPES.SMALL_PACKAGE && (
                      <ShipmentCustomerSIT
                        sitEstimatedWeight={mtoShipment.ppmShipment?.sitEstimatedWeight}
                        sitEstimatedEntryDate={mtoShipment.ppmShipment?.sitEstimatedEntryDate}
                        sitEstimatedDepartureDate={mtoShipment.ppmShipment?.sitEstimatedDepartureDate}
                      />
                    )}
                    <ShipmentWeight
                      authorizedWeight={serviceMember.weightAllotment.totalWeightSelf.toString()}
                      onEstimatedWeightChange={updateEstimatedWeightValue}
                    />
                  </>
                )}

                {isPPM && isAdvancePage && isServiceCounselor && mtoShipment.ppmShipment?.sitExpected && (
                  <SITCostDetails
                    cost={mtoShipment.ppmShipment?.sitEstimatedCost}
                    weight={mtoShipment.ppmShipment?.sitEstimatedWeight}
                    sitLocation={mtoShipment.ppmShipment?.sitLocation}
                    originZip={mtoShipment.ppmShipment?.pickupAddress.postalCode}
                    destinationZip={mtoShipment.ppmShipment?.destinationAddress.postalCode}
                    departureDate={mtoShipment.ppmShipment?.sitEstimatedDepartureDate}
                    entryDate={mtoShipment.ppmShipment?.sitEstimatedEntryDate}
                  />
                )}

                {isPPM && isAdvancePage && (
                  <ShipmentIncentiveAdvance
                    values={values}
                    estimatedIncentive={mtoShipment.ppmShipment?.estimatedIncentive}
                    advanceAmountRequested={mtoShipment.ppmShipment?.advanceAmountRequested}
                  />
                )}

                {(!isPPM || (isPPM && isAdvancePage)) && (
                  <ShipmentFormRemarks
                    userRole={userRole}
                    shipmentType={shipmentType}
                    customerRemarks={mtoShipment.customerRemarks}
                    counselorRemarks={mtoShipment.counselorRemarks}
                    showHint={false}
                    error={
                      formikProps.errors.counselorRemarks &&
                      (values.advanceRequested !== mtoShipment.ppmShipment?.hasRequestedAdvance ||
                        values.advance !== mtoShipment.ppmShipment?.advanceAmountRequested)
                    }
                  />
                )}

                {showAccountingCodes && (
                  <ShipmentAccountingCodes
                    TACs={TACs}
                    SACs={SACs}
                    onEditCodesClick={() => navigate(editOrdersPath)}
                    optional={isServiceCounselor}
                  />
                )}

                <div className={`${formStyles.formActions} ${styles.buttonGroup}`}>
                  {!isPPM && (
                    <Button
                      data-testid="submitForm"
                      disabled={isSubmitting || !isValid}
                      type="submit"
                      onClick={handleSubmit}
                    >
                      Save
                    </Button>
                  )}
                  <Button
                    type="button"
                    secondary
                    onClick={() => {
                      navigate(moveDetailsPath);
                    }}
                  >
                    Cancel
                  </Button>
                  {isPPM && (
                    <Button
                      data-testid="submitForm"
                      disabled={isSubmitting || !isValid}
                      type="submit"
                      onClick={handleSubmit}
                    >
                      Save and Continue
                    </Button>
                  )}
                </div>
              </Form>
            </div>
          </>
        );
      }}
    </Formik>
  );
};

ShipmentForm.propTypes = {
  submitHandler: func.isRequired,
  onUpdate: func,
  isCreatePage: bool,
  isForServicesCounseling: bool,
  currentResidence: AddressShape.isRequired,
  newDutyLocationAddress: SimpleAddressShape,
  shipmentType: string.isRequired,
  mtoShipment: ShipmentShape,
  moveTaskOrderID: string.isRequired,
  mtoShipments: arrayOf(ShipmentShape).isRequired,
  serviceMember: shape({
    weightAllotment: shape({
      totalWeightSelf: number,
      ubAllowance: number,
    }),
    agency: string.isRequired,
  }).isRequired,
  TACs: AccountingCodesShape,
  SACs: AccountingCodesShape,
  userRole: oneOf(officeRoles).isRequired,
  displayDestinationType: bool,
  isAdvancePage: bool,
  move: shape({
    eTag: string,
    id: string,
    closeoutOffice: TransportationOfficeShape,
  }),
};

ShipmentForm.defaultProps = {
  isCreatePage: false,
  isForServicesCounseling: false,
  onUpdate: () => {},
  newDutyLocationAddress: {
    city: '',
    state: '',
    postalCode: '',
  },
  mtoShipment: ShipmentShape,
  TACs: {},
  SACs: {},
  displayDestinationType: false,
  isAdvancePage: false,
  move: {},
};

export default ShipmentForm;
