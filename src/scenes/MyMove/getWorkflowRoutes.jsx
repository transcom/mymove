import React from 'react';
import { Route } from 'react-router-dom';
import { every, some, get, findKey, pick } from 'lodash';
import ValidatedPrivateRoute from 'shared/User/ValidatedPrivateRoute';
import WizardPage from 'shared/WizardPage';
import generatePath from 'shared/WizardPage/generatePath';
import { no_op } from 'shared/utils';
import { NULL_UUID, SHIPMENT_OPTIONS, CONUS_STATUS } from 'shared/constants';
import DodInfo from 'scenes/ServiceMembers/DodInfo';
import SMName from 'scenes/ServiceMembers/Name';
import ContactInfo from 'scenes/ServiceMembers/ContactInfo';
import ResidentialAddress from 'scenes/ServiceMembers/ResidentialAddress';
import BackupMailingAddress from 'scenes/ServiceMembers/BackupMailingAddress';
import BackupContact from 'scenes/ServiceMembers/BackupContact';
import ProfileReview from 'scenes/Review/ProfileReview';

import Orders from 'scenes/Orders/Orders';
import DutyStation from 'scenes/ServiceMembers/DutyStation';

import TransitionToMove from 'scenes/Orders/TransitionToMove';
import UploadOrders from 'scenes/Orders/UploadOrders';

import MoveLanding from 'pages/MyMove/MoveLanding';
import SelectMoveType from 'pages/MyMove/SelectMoveType';
import ConusOrNot from 'pages/MyMove/ConusOrNot';
import MovingInfo from 'pages/MyMove/MovingInfo';

import PpmDateAndLocations from 'scenes/Moves/Ppm/DateAndLocation';
import PpmWeight from 'scenes/Moves/Ppm/Weight';
import Review from 'scenes/Review/Review';
import Agreement from 'scenes/Legalese';

import HHGMoveSetup from 'pages/MyMove/HHGMoveSetup';

const PageNotInFlow = ({ location }) => (
  <div className="usa-grid">
    <h3>Missing Context</h3>
    You are trying to load a page that the system does not have context for. Please go to the home page and try again.
  </div>
);

// USE THESE FOR STUBBING OUT FUTURE WORK
// const Placeholder = props => {
//   return (
//     <WizardPage
//       handleSubmit={() => undefined}
//       pageList={props.pageList}
//       pageKey={props.pageKey}
//     >
//       <div className="Todo-phase2">
//         <h1>Placeholder for {props.title}</h1>
//         <h2>{props.description}</h2>
//       </div>
//     </WizardPage>
//   );
// };

// const stub = (key, pages, description) => ({ match }) => (
//   <Placeholder
//     pageList={pages}
//     pageKey={key}
//     title={key}
//     description={description}
//   />
// );

const always = () => true;
// Todo: update this when moves can be completed
const myFirstRodeo = (props) => !props.lastMoveIsCanceled;
const notMyFirstRodeo = (props) => props.lastMoveIsCanceled;
const hasPPM = ({ selectedMoveType }) => selectedMoveType !== null && selectedMoveType === SHIPMENT_OPTIONS.PPM;
const inHhgFlow = (props) => props.context.flags.hhgFlow;
const inGhcFlow = (props) => props.context.flags.ghcFlow;
const removeForDemo = (props) => props.context.flags.disableForDemo;
const isCurrentMoveSubmitted = ({ move }) => {
  return get(move, 'status', 'DRAFT') === 'SUBMITTED';
};

const pages = {
  '/service-member/:serviceMemberId/conus-status': {
    isInFlow: inGhcFlow,
    isComplete: ({ sm }) => sm.is_profile_complete || every([sm.rank, sm.edipi, sm.affiliation]),
    render: (key, pages, description, props) => ({ match }) => {
      return (
        <WizardPage
          handleSubmit={no_op}
          pageList={pages}
          pageKey={key}
          match={match}
          canMoveNext={props.conusStatus === CONUS_STATUS.CONUS}
        >
          <ConusOrNot conusStatus={props.conusStatus} />
        </WizardPage>
      );
    },
  },
  '/service-member/:serviceMemberId/create': {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || every([sm.rank, sm.edipi, sm.affiliation]),
    render: (key, pages) => ({ match }) => <DodInfo pages={pages} pageKey={key} match={match} />,
  },
  '/service-member/:serviceMemberId/name': {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || every([sm.first_name, sm.last_name]),
    render: (key, pages) => ({ match }) => <SMName pages={pages} pageKey={key} match={match} />,
  },
  '/service-member/:serviceMemberId/contact-info': {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) =>
      sm.is_profile_complete ||
      (every([sm.telephone, sm.personal_email]) && some([sm.phone_is_preferred, sm.email_is_preferred])),
    render: (key, pages) => ({ match }) => <ContactInfo pages={pages} pageKey={key} match={match} />,
  },
  '/service-member/:serviceMemberId/duty-station': {
    isInFlow: myFirstRodeo,

    // api for duty station always returns an object, even when duty station is not set
    // if there is no duty station, that object will have a null uuid
    isComplete: ({ sm }) => sm.is_profile_complete || get(sm, 'current_station.id', NULL_UUID) !== NULL_UUID,
    render: (key, pages) => ({ match }) => <DutyStation pages={pages} pageKey={key} match={match} />,
    description: 'current duty station',
  },
  '/service-member/:serviceMemberId/residence-address': {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || Boolean(sm.residential_address),
    render: (key, pages) => ({ match }) => <ResidentialAddress pages={pages} pageKey={key} match={match} />,
  },
  '/service-member/:serviceMemberId/backup-mailing-address': {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || Boolean(sm.backup_mailing_address),
    render: (key, pages) => ({ match }) => <BackupMailingAddress pages={pages} pageKey={key} match={match} />,
  },
  '/service-member/:serviceMemberId/backup-contacts': {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm, orders, move, ppm, backupContacts }) => {
      return sm.is_profile_complete || backupContacts.length > 0;
    },
    render: (key, pages) => ({ match }) => <BackupContact pages={pages} pageKey={key} match={match} />,
    description: 'Backup contacts',
  },
  '/service-member/:serviceMemberId/move-landing': {
    isInFlow: (props) => myFirstRodeo(props) && inGhcFlow(props) && !removeForDemo(props),
    isComplete: always,
    render: (key, pages) => () => {
      return (
        <WizardPage handleSubmit={no_op} pageList={pages} pageKey={key}>
          <MoveLanding />
        </WizardPage>
      );
    },
  },
  '/profile-review': {
    isInFlow: notMyFirstRodeo,
    isComplete: always,
    render: (key, pages) => ({ match }) => <ProfileReview pages={pages} pageKey={key} match={match} />,
  },
  '/orders/': {
    isInFlow: always,
    isComplete: ({ sm, orders }) =>
      every([
        orders.orders_type,
        orders.issue_date,
        orders.report_by_date,
        get(orders, 'new_duty_station.id', NULL_UUID) !== NULL_UUID,
      ]),
    render: (key, pages) => ({ match }) => <Orders pages={pages} pageKey={key} match={match} />,
  },
  '/orders/upload': {
    isInFlow: always,
    isComplete: ({ sm, orders, uploads }) =>
      get(orders, 'uploaded_orders.uploads', []).length > 0 || uploads.length > 0,
    render: (key, pages) => ({ match }) => <UploadOrders pages={pages} pageKey={key} match={match} />,
    description: 'Upload your orders',
  },
  '/orders/transition': {
    isInFlow: always,
    isComplete: always,
    render: (key, pages, description, props) => ({ match }) => {
      return (
        <WizardPage handleSubmit={no_op} pageList={pages} pageKey={key} additionalParams={{ moveId: props.moveId }}>
          <TransitionToMove />
        </WizardPage>
      );
    },
  },
  '/moves/:moveId/moving-info': {
    isInFlow: (props) => inHhgFlow(props) && !removeForDemo(props),
    isComplete: always,
    render: (key, pages) => () => {
      return (
        <WizardPage handleSubmit={no_op} pageList={pages} pageKey={key}>
          <MovingInfo />
        </WizardPage>
      );
    },
  },
  '/moves/:moveId/select-type': {
    // TODO: prevent user from hard-coding URL if they have a PPM or HHG existent?
    isInFlow: inHhgFlow,
    isComplete: ({ sm, orders, move }) => get(move, 'selected_move_type', null),
    render: (key, pages, props) => ({ match, history }) => (
      <SelectMoveType pageList={pages} pageKey={key} match={match} push={history.push} />
    ),
  },
  '/moves/:moveId/ppm-start': {
    isInFlow: (state) => {
      return state.selectedMoveType === SHIPMENT_OPTIONS.PPM;
    },
    isComplete: ({ sm, orders, move, ppm }) => {
      return ppm && every([ppm.original_move_date, ppm.pickup_postal_code, ppm.destination_postal_code]);
    },
    render: (key, pages) => ({ match }) => <PpmDateAndLocations pages={pages} pageKey={key} match={match} />,
  },
  '/moves/:moveId/ppm-incentive': {
    isInFlow: hasPPM,
    isComplete: ({ sm, orders, move, ppm }) =>
      get(ppm, 'weight_estimate', null) && get(ppm, 'weight_estimate', 0) !== 0,
    render: (key, pages) => ({ match }) => <PpmWeight pages={pages} pageKey={key} match={match} />,
  },
  '/moves/:moveId/hhg-start': {
    isInFlow: (state) => inHhgFlow && state.selectedMoveType === SHIPMENT_OPTIONS.HHG,
    isComplete: ({ sm, orders, move, ppm, mtoShipment }) => {
      return (
        mtoShipment &&
        every([
          mtoShipment.requestedPickupDate,
          mtoShipment.requestedDeliveryDate,
          mtoShipment.pickupAddress,
          mtoShipment.shipmentType,
        ])
      );
    },
    render: (key, pages, description, props) => ({ match, history }) => (
      <HHGMoveSetup pageList={pages} pageKey={key} match={match} push={history.push} />
    ),
  },
  '/moves/:moveId/review': {
    isInFlow: always,
    isComplete: ({ sm, orders, move, ppm, mtoShipment }) => isCurrentMoveSubmitted(move),
    render: (key, pages) => ({ match }) => <Review pages={pages} pageKey={key} match={match} />,
  },
  '/moves/:moveId/agreement': {
    isInFlow: always,
    isComplete: ({ sm, orders, move, ppm, mtoShipment }) => isCurrentMoveSubmitted(move),
    render: (key, pages, description, props) => ({ match }) => {
      return <Agreement pages={pages} pageKey={key} match={match} selectedMoveType={props.selectedMoveType} />;
    },
  },
};

export const getPagesInFlow = ({ selectedMoveType, conusStatus, lastMoveIsCanceled, context }) =>
  Object.keys(pages).filter((pageKey) => {
    // eslint-disable-next-line security/detect-object-injection
    const page = pages[pageKey];
    return page.isInFlow({ selectedMoveType, conusStatus, lastMoveIsCanceled, context });
  });

export const getNextIncompletePage = ({
  selectedMoveType = undefined,
  conusStatus = '',
  lastMoveIsCanceled = false,
  serviceMember = {},
  orders = {},
  uploads = [],
  move = {},
  ppm = {},
  mtoShipment = {},
  backupContacts = [],
  context = {},
}) => {
  const rawPath = findKey(
    pages,
    (p) =>
      p.isInFlow({ selectedMoveType, conusStatus, lastMoveIsCanceled, context }) &&
      !p.isComplete({ sm: serviceMember, orders, uploads, move, ppm, mtoShipment, backupContacts }),
  );
  const compiledPath = generatePath(rawPath, {
    serviceMemberId: get(serviceMember, 'id'),
    moveId: get(move, 'id'),
  });
  return compiledPath;
};

export const getWorkflowRoutes = (props) => {
  const flowProps = pick(props, ['selectedMoveType', 'conusStatus', 'lastMoveIsCanceled', 'context']);
  const pageList = getPagesInFlow(flowProps);
  return Object.keys(pages).map((key) => {
    // eslint-disable-next-line security/detect-object-injection
    const currPage = pages[key];
    if (currPage.isInFlow(flowProps)) {
      const render = currPage.render(key, pageList, currPage.description, props);
      return <ValidatedPrivateRoute exact path={key} key={key} render={render} />;
    } else {
      return <Route exact path={key} key={key} component={PageNotInFlow} />;
    }
  });
};
