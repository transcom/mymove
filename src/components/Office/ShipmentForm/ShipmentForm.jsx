import React, { useEffect, useState } from 'react';
import { arrayOf, bool, func, number, shape, string, oneOf } from 'prop-types';
import { Field, Formik } from 'formik';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { Alert, Button, Checkbox, Fieldset, FormGroup, Radio } from '@trussworks/react-uswds';
import classNames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

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
import { Form } from 'components/form/Form';
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
import { SHIPMENT_OPTIONS } from 'shared/constants';
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
import { datePickerFormat, formatDate } from 'shared/dates';

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
  const [successMessage, setSuccessMessage] = useState(null);
  const [shipmentAddressUpdateReviewErrorMessage, setShipmentAddressUpdateReviewErrorMessage] = useState(null);

  const [isCancelModalVisible, setIsCancelModalVisible] = useState(false);
  const [isAddressChangeModalOpen, setIsAddressChangeModalOpen] = useState(false);

  const [isTertiaryAddressEnabled, setIsTertiaryAddressEnabled] = useState(false);
  useEffect(() => {
    const fetchData = async () => {
      setIsTertiaryAddressEnabled(await isBooleanFlagEnabled('third_address_available'));
    };
    fetchData();
  }, []);

  const shipments = mtoShipments;

  const [isRequestedPickupDateAlertVisible, setIsRequestedPickupDateAlertVisible] = useState(false);
  const [isRequestedDeliveryDateAlertVisible, setIsRequestedDeliveryDateAlertVisible] = useState(false);
  const [requestedPickupDateAlertMessage, setRequestedPickupDateAlertMessage] = useState('');
  const [requestedDeliveryDateAlertMessage, setRequestedDeliveryDateAlertMessage] = useState('');
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
    // TODO - this is not supported by IE11, shipment number should be calculable from Redux anyways
    // we should fix this also b/c it doesn't display correctly in storybook
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
    shipmentType === SHIPMENT_OPTIONS.BOAT_HAUL_AWAY ||
    shipmentType === SHIPMENT_OPTIONS.BOAT_TOW_AWAY;

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
    if (isPPM) {
      const ppmShipmentBody = formatPpmShipmentForAPI(formValues);

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
              if (formValues.closeoutOffice.id) {
                mutateMoveCloseoutOffice(
                  {
                    locator: moveCode,
                    ifMatchETag: move.eTag,
                    body: { closeoutOfficeId: formValues.closeoutOffice.id },
                  },
                  {
                    onSuccess: () => {
                      actions.setSubmitting(false);
                      navigate(currentPath, { replace: true });
                      if (isTOO) {
                        navigate(moveViewPath);
                      } else {
                        navigate(advancePath);
                      }
                      setErrorMessage(null);
                      onUpdate('success');
                    },
                    onError: () => {
                      actions.setSubmitting(false);
                      setErrorMessage(`Something went wrong, and your changes were not saved. Please try again.`);
                    },
                  },
                );
              } else {
                navigate(currentPath, { replace: true });
                if (isTOO) {
                  navigate(moveViewPath);
                } else {
                  navigate(advancePath);
                }
              }
            },
            onError: () => {
              actions.setSubmitting(false);
              setErrorMessage(`Something went wrong, and your changes were not saved. Please try again.`);
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
                onError: () => {
                  actions.setSubmitting(false);
                  setErrorMessage(`Something went wrong, and your changes were not saved. Please try again.`);
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
        onError: () => {
          actions.setSubmitting(false);
          setErrorMessage(`Something went wrong, and your changes were not saved. Please try again.`);
        },
      });
      return;
    }

    //* MTO Shipments *//

    const {
      pickup,
      hasDeliveryAddress,
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

    const deliveryDetails = delivery;
    if (hasDeliveryAddress === 'no' && shipmentType !== SHIPMENT_OPTIONS.NTSR) {
      delete deliveryDetails.address;
    }

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
      delivery: deliveryDetails,
      ntsRecordedWeight,
      tacType: nullableTacType,
      sacType: nullableSacType,
      serviceOrderNumber,
      storageFacility,
      usesExternalVendor,
      destinationType,
      hasSecondaryPickup: hasSecondaryPickup === 'yes',
      secondaryPickup: hasSecondaryPickup === 'yes' ? secondaryPickup : {},
      hasSecondaryDelivery: hasSecondaryDelivery === 'yes',
      secondaryDelivery: hasSecondaryDelivery === 'yes' ? secondaryDelivery : {},
      hasTertiaryPickup: hasTertiaryPickup === 'yes',
      tertiaryPickup: hasTertiaryPickup === 'yes' ? tertiaryPickup : {},
      hasTertiaryDelivery: hasTertiaryDelivery === 'yes',
      tertiaryDelivery: hasTertiaryDelivery === 'yes' ? tertiaryDelivery : {},
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
          onError: () => {
            setErrorMessage(`Something went wrong, and your changes were not saved. Please try again.`);
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
        onError: () => {
          setErrorMessage(`Something went wrong, and your changes were not saved. Please try again.`);
        },
      });
    }
    // Edit a MTO Shipment as TOO
    else {
      submitHandler(updateMTOShipmentPayload, {
        onSuccess: () => {
          navigate(moveDetailsPath);
        },
        onError: () => {
          setErrorMessage(`Something went wrong, and your changes were not saved. Please try again.`);
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
      {({
        values,
        isValid,
        isSubmitting,
        setValues,
        handleSubmit,
        errors,
        touched,
        setFieldTouched,
        setFieldError,
        validateForm,
      }) => {
        const {
          hasSecondaryDestination,
          hasTertiaryDestination,
          hasDeliveryAddress,
          hasSecondaryPickup,
          hasSecondaryDelivery,
          hasTertiaryPickup,
          hasTertiaryDelivery,
        } = values;

        const lengthHasError = !!(
          (touched.lengthFeet && errors.lengthFeet === 'Required') ||
          (touched.lengthInches && errors.lengthFeet === 'Required')
        );
        const widthHasError = !!(
          (touched.widthFeet && errors.widthFeet === 'Required') ||
          (touched.widthInches && errors.widthFeet === 'Required')
        );
        const heightHasError = !!(
          (touched.heightFeet && errors.heightFeet === 'Required') ||
          (touched.heightInches && errors.heightFeet === 'Required')
        );
        const dimensionError = !!(
          (touched.lengthFeet && errors.lengthFeet?.includes('Dimensions')) ||
          (touched.lengthInches && errors.lengthFeet?.includes('Dimensions'))
        );
        if (touched.lengthInches && !touched.lengthFeet) {
          setFieldTouched('lengthFeet', true);
        }
        if (touched.widthInches && !touched.widthFeet) {
          setFieldTouched('widthFeet', true);
        }
        if (touched.heightInches && !touched.heightFeet) {
          setFieldTouched('heightFeet', true);
        }
        // manually turn off 'required' error when page loads if field is empty.
        if (values.year === null && !touched.year && errors.year === 'Required') {
          setFieldError('year', null);
        }

        const handleUseCurrentResidenceChange = (e) => {
          const { checked } = e.target;
          if (checked) {
            // use current residence
            setValues({
              ...values,
              pickup: {
                ...values.pickup,
                address: currentResidence,
              },
            });
          } else {
            // Revert address
            setValues({
              ...values,
              pickup: {
                ...values.pickup,
                address: {
                  streetAddress1: '',
                  streetAddress2: '',
                  streetAddress3: '',
                  city: '',
                  state: '',
                  postalCode: '',
                },
              },
            });
          }
        };

        const handlePickupDateChange = (e) => {
          setValues({
            ...values,
            pickup: {
              ...values.pickup,
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
            'Requested pickup date',
            setRequestedPickupDateAlertMessage,
            setIsRequestedPickupDateAlertVisible,
            onErrorHandler,
          );
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
            setErrorMessage(response?.body?.detail);
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
                {errorMessage}
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
                Request needs review. <a href="#delivery-location">See delivery location to proceed.</a>
              </Alert>
            )}

            <div className={styles.ShipmentForm}>
              <div className={styles.headerWrapper}>
                <div>
                  <ShipmentTag shipmentType={shipmentType} shipmentNumber={shipmentNumber} />

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

              <SectionWrapper className={styles.weightAllowance}>
                <p>
                  <strong>Weight allowance: </strong>
                  {formatWeight(serviceMember.weightAllotment.totalWeightSelf)}
                </p>
              </SectionWrapper>

              <Form className={formStyles.form}>
                {isTOO && !isHHG && !isPPM && !isBoat && !isMobileHome && <ShipmentVendor />}

                {isNTSR && <ShipmentWeightInput userRole={userRole} />}

                {isMobileHome && (
                  <MobileHomeShipmentForm
                    lengthHasError={lengthHasError}
                    widthHasError={widthHasError}
                    heightHasError={heightHasError}
                    values={values}
                    setFieldTouched={setFieldTouched}
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
                    setFieldTouched={setFieldTouched}
                    setFieldError={setFieldError}
                    validateForm={validateForm}
                    dimensionError={dimensionError}
                  />
                )}

                {showPickupFields && (
                  <SectionWrapper className={formStyles.formSection}>
                    <h2 className={styles.SectionHeaderExtraSpacing}>Pickup details</h2>
                    <Fieldset>
                      {isRequestedPickupDateAlertVisible && (
                        <Alert type="warning" aria-live="polite" headingLevel="h4">
                          {requestedPickupDateAlertMessage}
                        </Alert>
                      )}
                      <DatePickerInput
                        name="pickup.requestedDate"
                        label="Requested pickup date"
                        id="requestedPickupDate"
                        validate={validateDate}
                        onChange={handlePickupDateChange}
                      />
                    </Fieldset>
                    {!isNTSR && (
                      <>
                        <AddressFields
                          name="pickup.address"
                          legend="Pickup location"
                          render={(fields) => (
                            <>
                              <p>What address are the movers picking up from?</p>
                              <Checkbox
                                data-testid="useCurrentResidence"
                                label="Use current address"
                                name="useCurrentResidence"
                                onChange={handleUseCurrentResidenceChange}
                                id="useCurrentResidenceCheckbox"
                              />
                              {fields}
                              <h4>Second pickup location</h4>
                              <FormGroup>
                                <p>Do you want movers to pick up any belongings from a second address?</p>
                                <div className={formStyles.radioGroup}>
                                  <Field
                                    as={Radio}
                                    id="has-secondary-pickup"
                                    data-testid="has-secondary-pickup"
                                    label="Yes"
                                    name="hasSecondaryPickup"
                                    value="yes"
                                    title="Yes, I have a second pickup location"
                                    checked={hasSecondaryPickup === 'yes'}
                                  />
                                  <Field
                                    as={Radio}
                                    id="no-secondary-pickup"
                                    data-testid="no-secondary-pickup"
                                    label="No"
                                    name="hasSecondaryPickup"
                                    value="no"
                                    title="No, I do not have a second pickup location"
                                    checked={hasSecondaryPickup !== 'yes'}
                                  />
                                </div>
                              </FormGroup>
                              {hasSecondaryPickup === 'yes' && (
                                <>
                                  <AddressFields name="secondaryPickup.address" />
                                  {isTertiaryAddressEnabled && (
                                    <>
                                      <h4>Third pickup location</h4>
                                      <FormGroup>
                                        <p>Do you want movers to pick up any belongings from a third address?</p>
                                        <div className={formStyles.radioGroup}>
                                          <Field
                                            as={Radio}
                                            id="has-tertiary-pickup"
                                            data-testid="has-tertiary-pickup"
                                            label="Yes"
                                            name="hasTertiaryPickup"
                                            value="yes"
                                            title="Yes, I have a third pickup location"
                                            checked={hasTertiaryPickup === 'yes'}
                                          />
                                          <Field
                                            as={Radio}
                                            id="no-tertiary-pickup"
                                            data-testid="no-tertiary-pickup"
                                            label="No"
                                            name="hasTertiaryPickup"
                                            value="no"
                                            title="No, I do not have a third pickup location"
                                            checked={hasTertiaryPickup !== 'yes'}
                                          />
                                        </div>
                                      </FormGroup>
                                      {hasTertiaryPickup === 'yes' && <AddressFields name="tertiaryPickup.address" />}
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
                    <StorageFacilityAddress />
                  </>
                )}

                {isServiceCounselor && isNTSR && (
                  <>
                    <StorageFacilityInfo userRole={userRole} />
                    <StorageFacilityAddress />
                  </>
                )}

                {showDeliveryFields && (
                  <SectionWrapper className={formStyles.formSection}>
                    <h2 className={styles.SectionHeaderExtraSpacing}>Delivery details</h2>
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
                              Pending delivery location change request needs review.{' '}
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
                          legend="Delivery location"
                          disabled={deliveryAddressUpdateRequested}
                          className={classNames('usa-legend', styles.mockLegend)}
                        >
                          <AddressFields
                            name="delivery.address"
                            render={(fields) => {
                              return fields;
                            }}
                          />
                          <h4>Second delivery location</h4>
                          <FormGroup>
                            <p>Do you want the movers to deliver any belongings to a second address?</p>
                            <div className={formStyles.radioGroup}>
                              <Field
                                as={Radio}
                                data-testid="has-secondary-delivery"
                                id="has-secondary-delivery"
                                label="Yes"
                                name="hasSecondaryDelivery"
                                value="yes"
                                title="Yes, I have a second destination location"
                                checked={hasSecondaryDelivery === 'yes'}
                              />
                              <Field
                                as={Radio}
                                data-testid="no-secondary-delivery"
                                id="no-secondary-delivery"
                                label="No"
                                name="hasSecondaryDelivery"
                                value="no"
                                title="No, I do not have a second destination location"
                                checked={hasSecondaryDelivery !== 'yes'}
                              />
                            </div>
                          </FormGroup>
                          {hasSecondaryDelivery === 'yes' && (
                            <>
                              <AddressFields name="secondaryDelivery.address" />
                              {isTertiaryAddressEnabled && (
                                <>
                                  <h4>Third delivery location</h4>
                                  <FormGroup>
                                    <p>Do you want the movers to deliver any belongings from a third address?</p>
                                    <div className={formStyles.radioGroup}>
                                      <Field
                                        as={Radio}
                                        id="has-tertiary-delivery"
                                        data-testid="has-tertiary-delivery"
                                        label="Yes"
                                        name="hasTertiaryDelivery"
                                        value="yes"
                                        title="Yes, I have a third delivery location"
                                        checked={hasTertiaryDelivery === 'yes'}
                                      />
                                      <Field
                                        as={Radio}
                                        id="no-tertiary-delivery"
                                        data-testid="no-tertiary-delivery"
                                        label="No"
                                        name="hasTertiaryDelivery"
                                        value="no"
                                        title="No, I do not have a third delivery location"
                                        checked={hasTertiaryDelivery !== 'yes'}
                                      />
                                    </div>
                                  </FormGroup>
                                  {hasTertiaryDelivery === 'yes' && <AddressFields name="tertiaryDelivery.address" />}
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
                          Delivery location
                        </p>
                        {deliveryAddressUpdateRequested && (
                          <Alert type="error" slim className={styles.deliveryAddressUpdateAlert}>
                            <span className={styles.deliveryAddressUpdateAlertContent}>
                              Pending delivery location change request needs review.{' '}
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
                          legend="Delivery location"
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
                                value="yes"
                                title="Yes, I know my delivery address"
                                checked={hasDeliveryAddress === 'yes'}
                              />
                              <Field
                                as={Radio}
                                id="no-delivery-address"
                                label="No"
                                name="hasDeliveryAddress"
                                value="no"
                                title="No, I do not know my delivery address"
                                checked={hasDeliveryAddress === 'no'}
                              />
                            </div>
                          </FormGroup>
                          {hasDeliveryAddress === 'yes' ? (
                            <AddressFields
                              name="delivery.address"
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
                                  <h4>Second delivery location</h4>
                                  <FormGroup>
                                    <p>Do you want the movers to deliver any belongings to a second address?</p>
                                    <div className={formStyles.radioGroup}>
                                      <Field
                                        as={Radio}
                                        data-testid="has-secondary-delivery"
                                        id="has-secondary-delivery"
                                        label="Yes"
                                        name="hasSecondaryDelivery"
                                        value="yes"
                                        title="Yes, I have a second destination location"
                                        checked={hasSecondaryDelivery === 'yes'}
                                      />
                                      <Field
                                        as={Radio}
                                        data-testid="no-secondary-delivery"
                                        id="no-secondary-delivery"
                                        label="No"
                                        name="hasSecondaryDelivery"
                                        value="no"
                                        title="No, I do not have a second destination location"
                                        checked={hasSecondaryDelivery !== 'yes'}
                                      />
                                    </div>
                                  </FormGroup>
                                  {hasSecondaryDelivery === 'yes' && (
                                    <>
                                      <AddressFields name="secondaryDelivery.address" />
                                      {isTertiaryAddressEnabled && (
                                        <>
                                          <h4>Third delivery location</h4>
                                          <FormGroup>
                                            <p>
                                              Do you want the movers to deliver any belongings from a third address?
                                            </p>
                                            <div className={formStyles.radioGroup}>
                                              <Field
                                                as={Radio}
                                                id="has-tertiary-delivery"
                                                data-testid="has-tertiary-delivery"
                                                label="Yes"
                                                name="hasTertiaryDelivery"
                                                value="yes"
                                                title="Yes, I have a third delivery location"
                                                checked={hasTertiaryDelivery === 'yes'}
                                              />
                                              <Field
                                                as={Radio}
                                                id="no-tertiary-delivery"
                                                data-testid="no-tertiary-delivery"
                                                label="No"
                                                name="hasTertiaryDelivery"
                                                value="no"
                                                title="No, I do not have a third delivery location"
                                                checked={hasTertiaryDelivery !== 'yes'}
                                              />
                                            </div>
                                          </FormGroup>
                                          {hasTertiaryDelivery === 'yes' && (
                                            <AddressFields name="tertiaryDelivery.address" />
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
                    <SectionWrapper className={classNames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                      <h2>Departure date</h2>
                      <DatePickerInput name="expectedDepartureDate" label="Planned Departure Date" />
                      <Hint className={ppmStyles.hint}>
                        Enter the first day you expect to move things. It&apos;s OK if the actual date is different. We
                        will ask for your actual departure date when you document and complete your PPM.
                      </Hint>
                    </SectionWrapper>
                    <SectionWrapper className={classNames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                      <AddressFields
                        name="pickup.address"
                        legend="Pickup Address"
                        render={(fields) => (
                          <>
                            <p>What address are you moving from?</p>
                            <Checkbox
                              data-testid="useCurrentResidence"
                              label="Use Current Address"
                              name="useCurrentResidence"
                              onChange={handleUseCurrentResidenceChange}
                              id="useCurrentResidenceCheckbox"
                            />
                            {fields}
                            <h4>Second pickup address</h4>
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
                                  title="Yes, there is a second pickup location"
                                  checked={hasSecondaryPickup === 'true'}
                                />
                                <Field
                                  as={Radio}
                                  id="no-secondary-pickup"
                                  data-testid="no-secondary-pickup"
                                  label="No"
                                  name="hasSecondaryPickup"
                                  value="false"
                                  title="No, there is not a second pickup location"
                                  checked={hasSecondaryPickup !== 'true'}
                                />
                              </div>
                            </FormGroup>
                            {hasSecondaryPickup === 'true' && (
                              <>
                                <AddressFields name="secondaryPickup.address" />
                                {isTertiaryAddressEnabled && (
                                  <>
                                    <h4>Third pickup address</h4>
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
                                          title="Yes, there is a third pickup location"
                                          checked={hasTertiaryPickup === 'true'}
                                        />
                                        <Field
                                          as={Radio}
                                          id="no-tertiary-pickup"
                                          data-testid="no-tertiary-pickup"
                                          label="No"
                                          name="hasTertiaryPickup"
                                          value="false"
                                          title="No, there is not a third pickup location"
                                          checked={hasTertiaryPickup !== 'true'}
                                        />
                                      </div>
                                    </FormGroup>
                                    {hasTertiaryPickup === 'true' && <AddressFields name="tertiaryPickup.address" />}
                                  </>
                                )}
                              </>
                            )}
                          </>
                        )}
                      />
                      <AddressFields
                        name="destination.address"
                        legend="Delivery Address"
                        render={(fields) => (
                          <>
                            {fields}
                            <h4>Second delivery address</h4>
                            <FormGroup>
                              <p>
                                Will you move any belongings to a second address? (Must be near the delivery address.
                                Subject to approval.)
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
                                />
                              </div>
                            </FormGroup>
                            {hasSecondaryDestination === 'true' && (
                              <>
                                <AddressFields name="secondaryDestination.address" />
                                {isTertiaryAddressEnabled && (
                                  <>
                                    <h4>Third delivery address</h4>
                                    <FormGroup>
                                      <p>
                                        Will you move any belongings to a third address? (Must be near the delivery
                                        address. Subject to approval.)
                                      </p>
                                      <div className={formStyles.radioGroup}>
                                        <Field
                                          as={Radio}
                                          id="has-tertiary-destination"
                                          data-testid="has-tertiary-destination"
                                          label="Yes"
                                          name="hasTertiaryDestination"
                                          value="true"
                                          title="Yes, I have a third delivery location"
                                          checked={hasTertiaryDestination === 'true'}
                                        />
                                        <Field
                                          as={Radio}
                                          id="no-tertiary-destination"
                                          data-testid="no-tertiary-destination"
                                          label="No"
                                          name="hasTertiaryDestination"
                                          value="false"
                                          title="No, I do not have a third delivery location"
                                          checked={hasTertiaryDestination !== 'true'}
                                        />
                                      </div>
                                    </FormGroup>
                                    {hasTertiaryDestination === 'true' && (
                                      <AddressFields name="tertiaryDestination.address" />
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
                        <h2>Closeout office</h2>
                        <CloseoutOfficeInput
                          hint="If there is more than one PPM for this move, the closeout office will be the same for all your PPMs."
                          name="closeoutOffice"
                          placeholder="Start typing a closeout location..."
                          label="Closeout location"
                          displayAddress
                        />
                      </SectionWrapper>
                    )}
                    <ShipmentCustomerSIT
                      sitEstimatedWeight={mtoShipment.ppmShipment?.sitEstimatedWeight}
                      sitEstimatedEntryDate={mtoShipment.ppmShipment?.sitEstimatedEntryDate}
                      sitEstimatedDepartureDate={mtoShipment.ppmShipment?.sitEstimatedDepartureDate}
                    />
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
                      errors.counselorRemarks &&
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
