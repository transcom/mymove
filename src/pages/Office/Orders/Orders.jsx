import React, { useEffect, useReducer, useCallback, useState } from 'react';
import { Link, useNavigate, useParams, useLocation, generatePath } from 'react-router-dom';
import { Button, ErrorMessage } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { connect } from 'react-redux';

import ordersFormValidationSchema from './ordersFormValidationSchema';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import formStyles from 'styles/form.module.scss';
import { milmoveLogger } from 'utils/milmoveLog';
import { getTacValid, getLoa, updateOrder, getResponseError, getPayGradeOptions } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { tooRoutes, tioRoutes } from 'constants/routes';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { LineOfAccountingDfasElementOrder } from 'types/lineOfAccounting';
import OrdersDetailForm from 'components/Office/OrdersDetailForm/OrdersDetailForm';
import { formatSwaggerDate, dropdownInputOptions, formatPayGradeOptions } from 'utils/formatters';
import { DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { ORDERS_TYPE_DETAILS_OPTIONS, ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { ORDERS } from 'constants/queryKeys';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { LOA_VALIDATION_ACTIONS, reducer as loaReducer, initialState as initialLoaState } from 'reducers/loaValidation';
import { TAC_VALIDATION_ACTIONS, reducer as tacReducer, initialState as initialTacState } from 'reducers/tacValidation';
import { FEATURE_FLAG_KEYS, LOA_TYPE, MOVE_DOCUMENT_TYPE } from 'shared/constants';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import DocumentViewerFileManager from 'components/DocumentViewerFileManager/DocumentViewerFileManager';
import { setShowLoadingSpinner as setShowLoadingSpinnerAction } from 'store/general/actions';
import { scrollToViewFormikError } from 'utils/validation';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import retryPageLoading from 'utils/retryPageLoading';

const deptIndicatorDropdownOptions = dropdownInputOptions(DEPARTMENT_INDICATOR_OPTIONS);
const ordersTypeDetailsDropdownOptions = dropdownInputOptions(ORDERS_TYPE_DETAILS_OPTIONS);

const Orders = ({ files, amendedDocumentId, updateAmendedDocument, onAddFile, setShowLoadingSpinner }) => {
  const navigate = useNavigate();
  const { moveCode } = useParams();
  const [tacValidationState, tacValidationDispatch] = useReducer(tacReducer, null, initialTacState);
  const [loaValidationState, loaValidationDispatch] = useReducer(loaReducer, null, initialLoaState);
  const [orderTypesOptions, setOrderTypesOptions] = useState(ORDERS_TYPE_OPTIONS);
  const [serverError, setServerError] = useState(null);

  const { move, orders, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const { state } = useLocation();
  const orderId = move?.ordersId;
  const documentId = orders[orderId]?.uploaded_order_id;
  const amendedOrderDocumentId = orders[orderId]?.uploadedAmendedOrderID || amendedDocumentId;
  const from = state?.from;

  const ordersDocuments = files[MOVE_DOCUMENT_TYPE.ORDERS];
  const amendedDocuments = files[MOVE_DOCUMENT_TYPE.AMENDMENTS];
  const hasOrdersDocuments = ordersDocuments?.length > 0;
  const handleClose = useCallback(() => {
    let redirectPath;
    if (from === 'paymentRequestDetails') {
      redirectPath = generatePath(tioRoutes.BASE_PAYMENT_REQUESTS_PATH, { moveCode });
    } else {
      redirectPath = generatePath(tooRoutes.BASE_MOVE_VIEW_PATH, { moveCode });
    }
    navigate(redirectPath);
  }, [navigate, moveCode, from]);
  const queryClient = useQueryClient();
  const { mutate: mutateOrders } = useMutation(updateOrder, {
    onSuccess: (data, variables) => {
      const updatedOrder = data.orders[variables.orderID];
      queryClient.setQueryData([ORDERS, variables.orderID], {
        orders: {
          [`${variables.orderID}`]: updatedOrder,
        },
      });
      queryClient.invalidateQueries(ORDERS);
      handleClose();
    },
    onError: (error) => {
      const message = getResponseError(
        error,
        'Something went wrong, and your changes were not saved. Please refresh the page and try again.',
      );
      setServerError(message);
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  const buildFullLineOfAccountingString = (loa) => {
    const dfasMap = LineOfAccountingDfasElementOrder.map((key) => {
      if (key === 'loaEndFyTx') {
        // Specific logic for DFAS element A3, the loaEndFyTx.
        // This is a combination of both the BgFyTx and EndFyTx
        // and if one are null, then typically we would resort to "XXXXXXXX"
        // but for this one we'll just leave it empty as this is not for the EDI858.
        if (loa.loaBgFyTx != null && loa.loaEndFyTx != null) {
          return `${loa.loaBgFyTx}${loa.loaEndFyTx}`;
        }
        if (loa.loaBgFyTx === null || loa.loaByFyTx === undefined) {
          // Catch the scenario of loaBgFyTx being null but loaEndFyTx not being null
          return '';
        }
      }
      return loa[key] || '';
    });
    let longLoa = dfasMap.join('*');
    // remove any number of spaces following an asterisk in a LOA string
    longLoa = longLoa.replace(/\* +/g, '*');
    // remove any number of spaces preceding an asterisk in a LOA string
    longLoa = longLoa.replace(/ +\*/g, '*');

    return longLoa;
  };

  const { mutate: validateLoa } = useMutation(getLoa, {
    onSuccess: (data, variables) => {
      const { loaType } = variables;
      // The server decides if this is a valid LOA or not
      const isValid = (data?.validHhgProgramCodeForLoa ?? false) && (data?.validLoaForTac ?? false);
      // Construct the long line of accounting string
      const longLineOfAccounting = data ? buildFullLineOfAccountingString(data) : '';
      loaValidationDispatch({
        type: LOA_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        payload: {
          loa: data,
          longLineOfAccounting,
          isValid,
          loaType,
        },
      });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  const handleHHGTacValidation = async (value) => {
    if (value && value.length === 4 && value !== tacValidationState[LOA_TYPE.HHG].tac) {
      const response = await getTacValid({ tac: value });
      tacValidationDispatch({
        type: TAC_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        loaType: LOA_TYPE.HHG,
        isValid: response.isValid,
        tac: value,
      });
    }
  };

  const handleHHGLoaValidation = (formikValues) => {
    // LOA is not a field that can be interacted with
    // Validation is based on the scope of the form
    const { tac, issueDate, departmentIndicator } = formikValues;
    if (tac && tac.length === 4 && departmentIndicator && issueDate) {
      // Only run validation if a 4 length TAC is present, and department and issue date are also present
      validateLoa({
        tacCode: tac,
        departmentIndicator,
        effectiveDate: formatSwaggerDate(issueDate),
        loaType: LOA_TYPE.HHG,
      });
    }
  };

  const handleNTSLoaValidation = (formikValues) => {
    // LOA is not a field that can be interacted with
    // Validation is based on the scope of the form
    const { ntsTac, departmentIndicator } = formikValues;
    if (ntsTac && ntsTac.length === 4 && departmentIndicator) {
      // Only run validation if a 4 length NTS TAC and department are present
      // The effective date for an NTS LOA should be either the approved_at date of the
      // move, or the current time of review (Post review it will save as the approved_at)
      const effectiveDate = move?.approvedAt || Date.now();
      validateLoa({
        tacCode: ntsTac,
        departmentIndicator,
        effectiveDate: formatSwaggerDate(effectiveDate),
        loaType: LOA_TYPE.NTS,
      });
    }
  };

  const handleNTSTacValidation = async (value) => {
    if (value && value.length === 4 && value !== tacValidationState[LOA_TYPE.NTS].tac) {
      const response = await getTacValid({ tac: value });
      tacValidationDispatch({
        type: TAC_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        loaType: LOA_TYPE.NTS,
        isValid: response.isValid,
        tac: value,
      });
    }
  };

  const order = Object.values(orders)?.[0];

  const [payGradeDropdownOptions, setPayGradeOptions] = useState([]);
  useEffect(() => {
    const fetchGradeOptions = async () => {
      setShowLoadingSpinner(true, 'Loading Pay Grade options');
      try {
        const fetchedRanks = await getPayGradeOptions(order.agency);
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
  }, [order.agency, setShowLoadingSpinner]);

  const { entitlement, uploadedAmendedOrderID, amendedOrdersAcknowledgedAt } = order;
  // TODO - passing in these fields so they don't get unset. Need to rework the endpoint.
  const {
    proGearWeight,
    proGearWeightSpouse,
    requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment,
    dependentsAuthorized,
  } = entitlement;

  useEffect(() => {
    if (isLoading || isError) {
      return;
    }

    const checkHHGTac = async () => {
      const response = await getTacValid({ tac: order.tac });
      tacValidationDispatch({
        type: TAC_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        loaType: LOA_TYPE.HHG,
        isValid: response.isValid,
        tac: order.tac,
      });
    };

    const checkHHGLoa = () => {
      // Only run validation if a 4 length TAC is present, and department and issue date are also present
      validateLoa({
        tacCode: order?.tac,
        departmentIndicator: order?.department_indicator,
        effectiveDate: formatSwaggerDate(order?.date_issued),
        loaType: LOA_TYPE.HHG,
      });
    };

    const checkNTSLoa = () => {
      // Only run validation if a 4 length NTS TAC and department are present
      // The effective date for an NTS LOA should be either the approved_at date of the
      // move, or the current time of review (Post review it will save as the approved_at)
      const effectiveDate = move?.approvedAt || Date.now();
      validateLoa({
        tacCode: order?.ntsTac,
        departmentIndicator: order?.department_indicator,
        effectiveDate: formatSwaggerDate(effectiveDate),
        loaType: LOA_TYPE.NTS,
      });
    };

    const checkNTSTac = async () => {
      const response = await getTacValid({ tac: order.ntsTac });
      tacValidationDispatch({
        type: TAC_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        loaType: LOA_TYPE.NTS,
        isValid: response.isValid,
        tac: order.ntsTac,
      });
    };

    if (order?.tac && order.tac.length === 4) {
      checkHHGTac();
    }
    if (order?.ntsTac && order.ntsTac.length === 4) {
      checkNTSTac();
    }
    if (order?.tac && order.tac.length === 4 && order?.department_indicator && order?.date_issued) {
      checkHHGLoa();
    }
    if (order?.ntsTac && order?.ntsTac.length === 4 && order?.department_indicator) {
      checkNTSLoa();
    }
  }, [
    order?.tac,
    order?.ntsTac,
    order?.date_issued,
    order?.department_indicator,
    move?.approvedAt,
    isLoading,
    isError,
    validateLoa,
  ]);

  useEffect(() => {
    const checkFeatureFlags = async () => {
      const isWoundedWarriorEnabled = await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.WOUNDED_WARRIOR_MOVE);
      const isBluebarkEnabled = await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.BLUEBARK_MOVE);

      setOrderTypesOptions((prevOptions) => {
        const options = { ...prevOptions };
        if (!isWoundedWarriorEnabled) {
          delete options.WOUNDED_WARRIOR;
        }

        if (!isBluebarkEnabled) {
          delete options.BLUEBARK;
        }
        return options;
      });
    };

    checkFeatureFlags();
  }, []);
  const ordersTypeDropdownOptions = dropdownInputOptions(orderTypesOptions);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = (values) => {
    const { originDutyLocation, newDutyLocation, ...fields } = values;
    const body = {
      ...fields,
      originDutyLocationId: values.originDutyLocation.id,
      newDutyLocationId: values.newDutyLocation.id,
      issueDate: formatSwaggerDate(values.issueDate),
      reportByDate: formatSwaggerDate(values.reportByDate),
      proGearWeight,
      proGearWeightSpouse,
      requiredMedicalEquipmentWeight,
      organizationalClothingAndIndividualEquipment,
      grade: values.payGrade,
    };

    mutateOrders({ orderID: orderId, ifMatchETag: order.eTag, body });
  };

  const tacWarningMsg =
    'This TAC does not appear in TGET, so it might not be valid. Make sure it matches what‘s on the orders before you continue.';
  const hhgLoaMissingWarningMsg =
    'Unable to find a LOA based on the provided details. Please ensure an orders issue date, department indicator, and TAC are present on this form.';
  const ntsLoaMissingWarningMsg =
    'Unable to find a LOA based on the provided details. Please ensure a department indicator and TAC are present on this form.';
  const loaInvalidWarningMsg = 'The LOA identified based on the provided details appears to be invalid.';

  const hasAmendedOrders = !!uploadedAmendedOrderID;

  const initialValues = {
    agency: order?.agency,
    originDutyLocation: order?.originDutyLocation,
    newDutyLocation: order?.destinationDutyLocation,
    issueDate: order?.date_issued,
    reportByDate: order?.report_by_date,
    departmentIndicator: order?.department_indicator,
    ordersNumber: order?.order_number || '',
    ordersType: order?.order_type,
    ordersTypeDetail: order?.order_type_detail,
    tac: order?.tac || '',
    sac: order?.sac || '',
    ntsTac: order?.ntsTac,
    ntsSac: order?.ntsSac,
    ordersAcknowledgement: !!amendedOrdersAcknowledgedAt,
    payGrade: order?.grade,
    dependentsAuthorized,
  };

  return (
    <div className={styles.sidebar}>
      <Formik
        initialValues={initialValues}
        validationSchema={ordersFormValidationSchema}
        onSubmit={onSubmit}
        validateOnChange
      >
        {(formik) => {
          // onBlur, if the value has 4 digits, run validator and show warning if invalid
          const hhgTacWarning = tacValidationState[LOA_TYPE.HHG].isValid ? '' : tacWarningMsg;
          const ntsTacWarning = tacValidationState[LOA_TYPE.NTS].isValid ? '' : tacWarningMsg;
          // Conditionally set the LOA warning message based on off if it is missing or just invalid
          const isHHGLoaMissing =
            loaValidationState[LOA_TYPE.HHG].loa === null || loaValidationState[LOA_TYPE.HHG].loa === undefined;
          const isNTSLoaMissing =
            loaValidationState[LOA_TYPE.NTS].loa === null || loaValidationState[LOA_TYPE.NTS].loa === undefined;
          let hhgLoaWarning = '';
          let ntsLoaWarning = '';
          // Making a nested ternary here goes against linter rules
          // The primary warning should be if it is missing, the other warning should be if it is invalid
          if (isHHGLoaMissing) {
            hhgLoaWarning = hhgLoaMissingWarningMsg;
          } else if (!loaValidationState[LOA_TYPE.HHG].isValid) {
            hhgLoaWarning = loaInvalidWarningMsg;
          }
          if (isNTSLoaMissing) {
            ntsLoaWarning = ntsLoaMissingWarningMsg;
          } else if (!loaValidationState[LOA_TYPE.NTS].isValid) {
            ntsLoaWarning = loaInvalidWarningMsg;
          }

          return (
            <form onSubmit={formik.handleSubmit}>
              <div className={styles.content}>
                <div className={styles.top}>
                  <Button
                    className={styles.closeButton}
                    data-testid="closeSidebar"
                    type="button"
                    onClick={handleClose}
                    unstyled
                  >
                    <FontAwesomeIcon icon="times" title="Close sidebar" aria-label="Close sidebar" />
                  </Button>
                  <h2 className={styles.header}>View Orders</h2>
                  <div>
                    {/* the prop "from" represents the page the Link navigates away from.
                        Passing it via state allows the new page to have access to what the previous page was
                        for purposes of redirecting back to the previous page on cancel
                        For documentation please see https://reactrouter.com/en/main/components/link#state
                     */}
                    <Link
                      className={styles.viewAllowances}
                      data-testid="view-allowances"
                      to="../allowances"
                      state={{ from }}
                    >
                      View Allowances
                    </Link>
                  </div>
                  <Restricted to={permissionTypes.updateOrders}>
                    <DocumentViewerFileManager
                      fileUploadRequired={!hasOrdersDocuments}
                      orderId={orderId}
                      documentId={documentId}
                      files={ordersDocuments}
                      documentType={MOVE_DOCUMENT_TYPE.ORDERS}
                      onAddFile={onAddFile}
                    />
                    <DocumentViewerFileManager
                      orderId={orderId}
                      documentId={amendedOrderDocumentId}
                      files={amendedDocuments}
                      documentType={MOVE_DOCUMENT_TYPE.AMENDMENTS}
                      updateAmendedDocument={updateAmendedDocument}
                      onAddFile={onAddFile}
                    />
                  </Restricted>
                </div>
                <div className={styles.body}>
                  <Restricted
                    to={permissionTypes.updateOrders}
                    fallback={
                      <OrdersDetailForm
                        deptIndicatorOptions={deptIndicatorDropdownOptions}
                        ordersTypeOptions={ordersTypeDropdownOptions}
                        ordersTypeDetailOptions={ordersTypeDetailsDropdownOptions}
                        hhgTacWarning={hhgTacWarning}
                        ntsTacWarning={ntsTacWarning}
                        validateHHGTac={handleHHGTacValidation}
                        validateNTSTac={handleNTSTacValidation}
                        hhgLoaWarning={hhgLoaWarning}
                        ntsLoaWarning={ntsLoaWarning}
                        validateHHGLoa={() => handleHHGLoaValidation(formik.values)}
                        validateNTSLoa={() => handleNTSLoaValidation(formik.values)}
                        hhgLongLineOfAccounting={loaValidationState[LOA_TYPE.HHG].longLineOfAccounting}
                        ntsLongLineOfAccounting={loaValidationState[LOA_TYPE.NTS].longLineOfAccounting}
                        showOrdersAcknowledgement={hasAmendedOrders}
                        ordersType={order.order_type}
                        setFieldValue={formik.setFieldValue}
                        payGradeOptions={payGradeDropdownOptions}
                        formIsDisabled
                      />
                    }
                  >
                    <OrdersDetailForm
                      deptIndicatorOptions={deptIndicatorDropdownOptions}
                      ordersTypeOptions={ordersTypeDropdownOptions}
                      ordersTypeDetailOptions={ordersTypeDetailsDropdownOptions}
                      hhgTacWarning={hhgTacWarning}
                      ntsTacWarning={ntsTacWarning}
                      validateHHGTac={handleHHGTacValidation}
                      validateNTSTac={handleNTSTacValidation}
                      hhgLoaWarning={hhgLoaWarning}
                      ntsLoaWarning={ntsLoaWarning}
                      validateHHGLoa={() => handleHHGLoaValidation(formik.values)}
                      validateNTSLoa={() => handleNTSLoaValidation(formik.values)}
                      hhgLongLineOfAccounting={loaValidationState[LOA_TYPE.HHG].longLineOfAccounting}
                      ntsLongLineOfAccounting={loaValidationState[LOA_TYPE.NTS].longLineOfAccounting}
                      showOrdersAcknowledgement={hasAmendedOrders}
                      ordersType={order.order_type}
                      setFieldValue={formik.setFieldValue}
                      payGradeOptions={payGradeDropdownOptions}
                    />
                  </Restricted>
                </div>
                {serverError && <ErrorMessage>{serverError}</ErrorMessage>}
                <Restricted to={permissionTypes.updateOrders}>
                  <div className={styles.bottom}>
                    <div className={formStyles.formActions}>
                      <Button type="button" secondary onClick={handleClose}>
                        Cancel
                      </Button>
                      <Button disabled={formik.isSubmitting} type="submit" onClick={scrollToViewFormikError(formik)}>
                        Save
                      </Button>
                    </div>
                  </div>
                </Restricted>
              </div>
            </form>
          );
        }}
      </Formik>
    </div>
  );
};

const mapDispatchToProps = {
  setShowLoadingSpinner: setShowLoadingSpinnerAction,
};

export default connect(() => {}, mapDispatchToProps)(Orders);
