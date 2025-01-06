import React, { Component, lazy } from 'react';
import PropTypes from 'prop-types';
import { Route, Routes, Navigate } from 'react-router-dom';
import { isBooleanFlagEnabled } from '../../utils/featureFlags';
import { connect } from 'react-redux';
import { GovBanner } from '@trussworks/react-uswds';

import 'styles/full_uswds.scss';
import 'styles/customer.scss';

import { getWorkflowRoutes } from './getWorkflowRoutes';

// Logger
import { milmoveLogger } from 'utils/milmoveLog';
import { retryPageLoading } from 'utils/retryPageLoading';
import BypassBlock from 'components/BypassBlock';
import CUIHeader from 'components/CUIHeader/CUIHeader';
import LoggedOutHeader from 'containers/Headers/LoggedOutHeader';
import CustomerLoggedInHeader from 'containers/Headers/CustomerLoggedInHeader';
import Alert from 'shared/Alert';
import Footer from 'components/Customer/Footer';
import ConnectedLogoutOnInactivity from 'layout/LogoutOnInactivity';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { loadInternalSchema } from 'shared/Swagger/ducks';
import { withContext } from 'shared/AppContext';
import { no_op } from 'shared/utils';
import { generatePageTitle } from 'hooks/custom';
import { loadUser as loadUserAction } from 'store/auth/actions';
import { initOnboarding as initOnboardingAction } from 'store/onboarding/actions';
import {
  selectCacValidated,
  selectGetCurrentUserIsLoading,
  selectIsLoggedIn,
  selectUnderMaintenance,
} from 'store/auth/selectors';
import { selectConusStatus } from 'store/onboarding/selectors';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentMove,
  selectHasCanceledMove,
} from 'store/entities/selectors';
import { generalRoutes, customerRoutes } from 'constants/routes';
import { pageNames } from 'constants/signInPageNames';
/** Pages */
import InfectedUpload from 'shared/Uploader/InfectedUpload';
import ProcessingUpload from 'shared/Uploader/ProcessingUpload';
import Edit from 'scenes/Review/Edit';
import NotFound from 'components/NotFound/NotFound';
import PrivacyPolicyStatement from 'components/Statements/PrivacyAndPolicyStatement';
import AccessibilityStatement from 'components/Statements/AccessibilityStatement';
import ConnectedCreateOrEditMtoShipment from 'pages/MyMove/CreateOrEditMtoShipment';
import Home from 'pages/MyMove/Home';
import TitleAnnouncer from 'components/TitleAnnouncer/TitleAnnouncer';
import MultiMovesLandingPage from 'pages/MyMove/Multi-Moves/MultiMovesLandingPage';
import MoveHome from 'pages/MyMove/Home/MoveHome';
import AddOrders from 'pages/MyMove/AddOrders';
import UploadOrders from 'pages/MyMove/UploadOrders';
import SmartCardRedirect from 'shared/SmartCardRedirect/SmartCardRedirect';
import OktaErrorBanner from 'components/OktaErrorBanner/OktaErrorBanner';
import MaintenancePage from 'pages/Maintenance/MaintenancePage';
// Pages should be lazy-loaded (they correspond to unique routes & only need to be loaded when that URL is accessed)
const SignIn = lazy(() => import('pages/SignIn/SignIn'));
const InvalidPermissions = lazy(() => import('pages/InvalidPermissions/InvalidPermissions'));
const MovingInfo = lazy(() => import('pages/MyMove/MovingInfo'));
const EditServiceInfo = lazy(() => import('pages/MyMove/Profile/EditServiceInfo'));
const Profile = lazy(() => import('pages/MyMove/Profile/Profile'));
const EditContactInfo = lazy(() => import('pages/MyMove/Profile/EditContactInfo'));
const EditOktaInfo = lazy(() => import('pages/MyMove/Profile/EditOktaInfo'));
const AmendOrders = lazy(() => import('pages/MyMove/AmendOrders/AmendOrders'));
const EditOrders = lazy(() => import('pages/MyMove/EditOrders'));
const BoatShipmentLocationInfo = lazy(() =>
  import('pages/MyMove/Boat/BoatShipmentLocationInfo/BoatShipmentLocationInfo'),
);
const MobileHomeShipmentLocationInfo = lazy(() =>
  import('pages/MyMove/MobileHome/MobileHomeShipmentLocationInfo/MobileHomeShipmentLocationInfo'),
);
const EstimatedWeightsProGear = lazy(() =>
  import('pages/MyMove/PPM/Booking/EstimatedWeightsProGear/EstimatedWeightsProGear'),
);
const EstimatedIncentive = lazy(() => import('pages/MyMove/PPM/Booking/EstimatedIncentive/EstimatedIncentive'));
const Advance = lazy(() => import('pages/MyMove/PPM/Booking/Advance/Advance'));
const About = lazy(() => import('pages/MyMove/PPM/Closeout/About/About'));
const WeightTickets = lazy(() => import('pages/MyMove/PPM/Closeout/WeightTickets/WeightTickets'));
const PPMReview = lazy(() => import('pages/MyMove/PPM/Closeout/Review/Review'));
const ProGear = lazy(() => import('pages/MyMove/PPM/Closeout/ProGear/ProGear.jsx'));
const Expenses = lazy(() => import('pages/MyMove/PPM/Closeout/Expenses/Expenses'));
const PPMFinalCloseout = lazy(() => import('pages/MyMove/PPM/Closeout/FinalCloseout/FinalCloseout'));
const AdditionalDocuments = lazy(() => import('pages/MyMove/AdditionalDocuments/AdditionalDocuments'));
const PPMFeedback = lazy(() => import('pages/MyMove/PPM/Closeout/Feedback/Feedback'));

export class CustomerApp extends Component {
  constructor(props) {
    super(props);

    this.state = {
      hasError: false,
      error: undefined,
      info: undefined,
      multiMoveFeatureFlag: false,
      cacValidatedFeatureFlag: false,
      validationCodeRequired: false,
      oktaErrorBanner: false,
    };
  }

  componentDidMount() {
    const { loadUser, initOnboarding, loadInternalSchema } = this.props;

    loadInternalSchema();
    loadUser();
    initOnboarding();
    isBooleanFlagEnabled('multi_move').then((enabled) => {
      this.setState({
        multiMoveFeatureFlag: enabled,
      });
    });
    isBooleanFlagEnabled('cac_validated_login').then((enabled) => {
      this.setState({
        cacValidatedFeatureFlag: enabled,
      });
    });
    isBooleanFlagEnabled('validation_code_required').then((enabled) => {
      this.setState({
        validationCodeRequired: enabled,
      });
    });
    // if the params "okta_error=true" are appended to the url, then we need to change state to display a banner
    // this occurs when a user is trying to use an office user's email to access the customer application
    // Okta config rules do not allow the same email to be used for both office & customer apps
    const currentUrl = new URL(window.location.href);
    const oktaErrorParam = currentUrl.searchParams.get('okta_error');
    if (oktaErrorParam === 'true') {
      this.setState({
        oktaErrorBanner: true,
      });
    }
    document.title = generatePageTitle('Sign In');
  }

  componentDidCatch(error, info) {
    const { message } = error;
    milmoveLogger.error({ message, info });
    this.setState({
      hasError: true,
      error,
      info,
    });
    retryPageLoading(error);
  }

  render() {
    const { props } = this;
    const { userIsLoggedIn, loginIsLoading, cacValidated, underMaintenance } = props;
    const { hasError, multiMoveFeatureFlag, cacValidatedFeatureFlag, oktaErrorBanner } = this.state;
    const script = document.createElement('script');

    script.src = '//rum-static.pingdom.net/pa-6567b05deff3250012000426.js';
    script.async = true;
    document.body.appendChild(script);

    if (underMaintenance) {
      return <MaintenancePage />;
    }

    return (
      <>
        <div className="my-move site" id="app-root">
          <TitleAnnouncer />
          <CUIHeader />
          <BypassBlock />
          <GovBanner />

          {userIsLoggedIn ? <CustomerLoggedInHeader /> : <LoggedOutHeader app={pageNames.MYMOVE} />}

          <main role="main" className="site__content my-move-container" id="main">
            <ConnectedLogoutOnInactivity />

            <div className="usa-grid">
              {props.swaggerError && (
                <div className="grid-container">
                  <div className="grid-row">
                    <div className="grid-col-12">
                      <Alert type="error" heading="An error occurred">
                        There was an error contacting the server.
                      </Alert>
                    </div>
                  </div>
                </div>
              )}
            </div>

            {oktaErrorBanner && <OktaErrorBanner />}

            {hasError && <SomethingWentWrong />}

            {/* Showing Smart Card info page until user signs in with SC one time */}
            {userIsLoggedIn && !cacValidated && cacValidatedFeatureFlag && <SmartCardRedirect />}

            {/* No Auth Routes */}
            {!userIsLoggedIn && (
              <Routes>
                <Route path={generalRoutes.SIGN_IN_PATH} element={<SignIn />} />
                <Route path={generalRoutes.PRIVACY_SECURITY_POLICY_PATH} element={<PrivacyPolicyStatement />} />
                <Route path={generalRoutes.ACCESSIBILITY_PATH} element={<AccessibilityStatement />} />
                <Route
                  end
                  path="/forbidden"
                  element={
                    <div className="usa-grid">
                      <h2>You are forbidden to use this endpoint</h2>
                    </div>
                  }
                />
                <Route
                  end
                  path="/server_error"
                  element={
                    <div className="usa-grid">
                      <h2>We are experiencing an internal server error</h2>
                    </div>
                  }
                />
                <Route end path="/invalid-permissions" element={<InvalidPermissions />} />
                <Route
                  path="*"
                  element={(loginIsLoading && <LoadingPlaceholder />) || <Navigate to="/sign-in" replace />}
                />
              </Routes>
            )}

            {/* when the cacValidated feature flag is on, we need to check for the cacValidated value for rendering */}
            {cacValidatedFeatureFlag
              ? !hasError &&
                !props.swaggerError &&
                userIsLoggedIn &&
                cacValidated && (
                  <Routes>
                    {/* no auth routes should still exist */}
                    <Route path={generalRoutes.SIGN_IN_PATH} element={<SignIn />} />
                    <Route path={generalRoutes.MULTI_MOVES_LANDING_PAGE} element={<MultiMovesLandingPage />} />
                    <Route path={generalRoutes.PRIVACY_SECURITY_POLICY_PATH} element={<PrivacyPolicyStatement />} />
                    <Route path={generalRoutes.ACCESSIBILITY_PATH} element={<AccessibilityStatement />} />

                    {/* auth required */}
                    {/* <Route end path="/ppm" element={<PpmLanding />} /> */}

                    {/* ROOT */}
                    {/* If multiMove is enabled home page will route to dashboard element. Otherwise, it will route to the move page. */}
                    {multiMoveFeatureFlag ? (
                      <Route path={generalRoutes.HOME_PATH} end element={<MultiMovesLandingPage />} />
                    ) : (
                      <Route path={generalRoutes.HOME_PATH} end element={<Home />} />
                    )}

                    {getWorkflowRoutes(props)}

                    <Route end path={customerRoutes.MOVE_HOME_PAGE} element={<Home />} />
                    <Route end path={customerRoutes.MOVE_HOME_PATH} element={<MoveHome />} />
                    <Route end path={customerRoutes.SHIPMENT_MOVING_INFO_PATH} element={<MovingInfo />} />
                    <Route end path="/moves/:moveId/edit" element={<Edit />} />
                    <Route end path={customerRoutes.SERVICE_INFO_EDIT_PATH} element={<EditServiceInfo />} />
                    <Route path={customerRoutes.SHIPMENT_CREATE_PATH} element={<ConnectedCreateOrEditMtoShipment />} />
                    <Route end path={customerRoutes.PROFILE_PATH} element={<Profile />} />
                    <Route
                      end
                      path={customerRoutes.SHIPMENT_EDIT_PATH}
                      element={<ConnectedCreateOrEditMtoShipment />}
                    />
                    <Route
                      path={customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH}
                      element={<EstimatedWeightsProGear />}
                    />
                    <Route
                      end
                      path={customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH}
                      element={<EstimatedIncentive />}
                    />
                    <Route end path={customerRoutes.SHIPMENT_PPM_ADVANCES_PATH} element={<Advance />} />
                    <Route end path={customerRoutes.CONTACT_INFO_EDIT_PATH} element={<EditContactInfo />} />
                    <Route end path={customerRoutes.EDIT_OKTA_PROFILE_PATH} element={<EditOktaInfo />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_ABOUT_PATH} element={<About />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH} element={<WeightTickets />} />
                    <Route
                      end
                      path={customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH}
                      element={<WeightTickets />}
                    />
                    <Route end path={customerRoutes.SHIPMENT_PPM_REVIEW_PATH} element={<PPMReview />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_EXPENSES_PATH} element={<Expenses />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH} element={<Expenses />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_COMPLETE_PATH} element={<PPMFinalCloseout />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_FEEDBACK_PATH} element={<PPMFeedback />} />
                    <Route path={customerRoutes.ORDERS_ADD_PATH} element={<AddOrders />} />
                    <Route path={customerRoutes.ORDERS_EDIT_PATH} element={<EditOrders />} />
                    <Route path={customerRoutes.ORDERS_UPLOAD_PATH} element={<UploadOrders />} />
                    <Route path={customerRoutes.ORDERS_AMEND_PATH} element={<AmendOrders />} />
                    <Route
                      end
                      path={customerRoutes.UPLOAD_ADDITIONAL_DOCUMENTS_PATH}
                      element={<AdditionalDocuments />}
                    />
                    <Route end path="/infected-upload" element={<InfectedUpload />} />
                    <Route end path="/processing-upload" element={<ProcessingUpload />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_PRO_GEAR_PATH} element={<ProGear />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH} element={<ProGear />} />

                    {/* Errors */}
                    <Route
                      end
                      path="/forbidden"
                      element={
                        <div className="usa-grid">
                          <h2>You are forbidden to use this endpoint</h2>
                        </div>
                      }
                    />
                    <Route
                      end
                      path="/server_error"
                      element={
                        <div className="usa-grid">
                          <h2>We are experiencing an internal server error</h2>
                        </div>
                      }
                    />
                    <Route end path="/invalid-permissions" element={<InvalidPermissions />} />

                    {/* 404 - user logged in but at unknown route */}
                    <Route path="*" element={<NotFound />} />
                  </Routes>
                )
              : !hasError &&
                !props.swaggerError &&
                userIsLoggedIn && (
                  <Routes>
                    {/* no auth routes should still exist */}
                    <Route path={generalRoutes.SIGN_IN_PATH} element={<SignIn />} />
                    <Route path={generalRoutes.MULTI_MOVES_LANDING_PAGE} element={<MultiMovesLandingPage />} />
                    <Route path={generalRoutes.PRIVACY_SECURITY_POLICY_PATH} element={<PrivacyPolicyStatement />} />
                    <Route path={generalRoutes.ACCESSIBILITY_PATH} element={<AccessibilityStatement />} />

                    {/* auth required */}
                    {/* <Route end path="/ppm" element={<PpmLanding />} /> */}

                    {/* ROOT */}
                    {/* If multiMove is enabled home page will route to dashboard element. Otherwise, it will route to the move page. */}
                    {multiMoveFeatureFlag ? (
                      <Route path={generalRoutes.HOME_PATH} end element={<MultiMovesLandingPage />} />
                    ) : (
                      <Route path={generalRoutes.HOME_PATH} end element={<Home />} />
                    )}

                    {getWorkflowRoutes(props)}

                    <Route end path={customerRoutes.MOVE_HOME_PAGE} element={<Home />} />
                    <Route end path={customerRoutes.MOVE_HOME_PATH} element={<MoveHome />} />
                    <Route end path={customerRoutes.SHIPMENT_MOVING_INFO_PATH} element={<MovingInfo />} />
                    <Route end path="/moves/:moveId/edit" element={<Edit />} />
                    <Route end path={customerRoutes.SERVICE_INFO_EDIT_PATH} element={<EditServiceInfo />} />
                    <Route path={customerRoutes.SHIPMENT_CREATE_PATH} element={<ConnectedCreateOrEditMtoShipment />} />
                    <Route end path={customerRoutes.PROFILE_PATH} element={<Profile />} />
                    <Route
                      end
                      path={customerRoutes.SHIPMENT_EDIT_PATH}
                      element={<ConnectedCreateOrEditMtoShipment />}
                    />
                    <Route path={customerRoutes.SHIPMENT_BOAT_LOCATION_INFO} element={<BoatShipmentLocationInfo />} />
                    <Route
                      path={customerRoutes.SHIPMENT_MOBILE_HOME_LOCATION_INFO}
                      element={<MobileHomeShipmentLocationInfo />}
                    />
                    <Route
                      path={customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH}
                      element={<EstimatedWeightsProGear />}
                    />
                    <Route
                      end
                      path={customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH}
                      element={<EstimatedIncentive />}
                    />
                    <Route end path={customerRoutes.SHIPMENT_PPM_ADVANCES_PATH} element={<Advance />} />
                    <Route end path={customerRoutes.CONTACT_INFO_EDIT_PATH} element={<EditContactInfo />} />
                    <Route end path={customerRoutes.EDIT_OKTA_PROFILE_PATH} element={<EditOktaInfo />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_ABOUT_PATH} element={<About />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH} element={<WeightTickets />} />
                    <Route
                      end
                      path={customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH}
                      element={<WeightTickets />}
                    />
                    <Route end path={customerRoutes.SHIPMENT_PPM_REVIEW_PATH} element={<PPMReview />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_EXPENSES_PATH} element={<Expenses />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH} element={<Expenses />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_COMPLETE_PATH} element={<PPMFinalCloseout />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_FEEDBACK_PATH} element={<PPMFeedback />} />
                    <Route path={customerRoutes.ORDERS_ADD_PATH} element={<AddOrders />} />
                    <Route path={customerRoutes.ORDERS_EDIT_PATH} element={<EditOrders />} />
                    <Route path={customerRoutes.ORDERS_UPLOAD_PATH} element={<UploadOrders />} />
                    <Route path={customerRoutes.ORDERS_AMEND_PATH} element={<AmendOrders />} />
                    <Route
                      end
                      path={customerRoutes.UPLOAD_ADDITIONAL_DOCUMENTS_PATH}
                      element={<AdditionalDocuments />}
                    />
                    <Route end path="/infected-upload" element={<InfectedUpload />} />
                    <Route end path="/processing-upload" element={<ProcessingUpload />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_PRO_GEAR_PATH} element={<ProGear />} />
                    <Route end path={customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH} element={<ProGear />} />

                    {/* Errors */}
                    <Route
                      end
                      path="/forbidden"
                      element={
                        <div className="usa-grid">
                          <h2>You are forbidden to use this endpoint</h2>
                        </div>
                      }
                    />
                    <Route
                      end
                      path="/server_error"
                      element={
                        <div className="usa-grid">
                          <h2>We are experiencing an internal server error</h2>
                        </div>
                      }
                    />
                    <Route end path="/invalid-permissions" element={<InvalidPermissions />} />

                    {/* 404 - user logged in but at unknown route */}
                    <Route path="*" element={<NotFound />} />
                  </Routes>
                )}
          </main>
          <Footer />
        </div>
        <div id="modal-root" />
      </>
    );
  }
}

CustomerApp.propTypes = {
  loadInternalSchema: PropTypes.func,
  loadUser: PropTypes.func,
  initOnboarding: PropTypes.func,
  loginIsLoading: PropTypes.bool,
  userIsLoggedIn: PropTypes.bool,
  conusStatus: PropTypes.string,
  context: PropTypes.shape({
    flags: PropTypes.shape({
      hhgFlow: PropTypes.bool,
      ghcFlow: PropTypes.bool,
    }),
  }).isRequired,
};

CustomerApp.defaultProps = {
  loadInternalSchema: no_op,
  loadUser: no_op,
  initOnboarding: no_op,
  loginIsLoading: false,
  userIsLoggedIn: false,
  conusStatus: '',
  context: {
    flags: {
      hhgFlow: false,
      ghcFlow: false,
    },
  },
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;
  const move = selectCurrentMove(state) || {};

  return {
    loginIsLoading: selectGetCurrentUserIsLoading(state),
    userIsLoggedIn: selectIsLoggedIn(state),
    cacValidated: selectCacValidated(serviceMember),
    currentServiceMemberId: serviceMemberId,
    lastMoveIsCanceled: selectHasCanceledMove(state),
    moveId: move?.id,
    conusStatus: selectConusStatus(state),
    swaggerError: state.swaggerInternal.hasErrored,
    underMaintenance: selectUnderMaintenance(state),
  };
};
const mapDispatchToProps = {
  loadInternalSchema,
  loadUser: loadUserAction,
  initOnboarding: initOnboardingAction,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(CustomerApp));
