import React from 'react';
import { Route } from 'react-router-dom';
import { every, some, get, findKey } from 'lodash';
import PrivateRoute from 'shared/User/PrivateRoute';
import WizardPage from 'shared/WizardPage';
import generatePath from 'shared/WizardPage/generatePath';
import { no_op } from 'shared/utils';

import DodInfo from 'scenes/ServiceMembers/DodInfo';
import SMName from 'scenes/ServiceMembers/Name';
import ContactInfo from 'scenes/ServiceMembers/ContactInfo';
import ResidentialAddress from 'scenes/ServiceMembers/ResidentialAddress';
import BackupMailingAddress from 'scenes/ServiceMembers/BackupMailingAddress';
import BackupContact from 'scenes/ServiceMembers/BackupContact';
import TransitionToOrders from 'scenes/ServiceMembers/TransitionToOrders';
import Orders from 'scenes/Orders/Orders';
import DutyStation from 'scenes/ServiceMembers/DutyStation';

import TransitionToMove from 'scenes/Orders/TransitionToMove';
import UploadOrders from 'scenes/Orders/UploadOrders';

import MoveType from 'scenes/Moves/MoveTypeWizard';
import Transition from 'scenes/Moves/Transition';
import PpmDateAndLocations from 'scenes/Moves/Ppm/DateAndLocation';
import PpmWeight from 'scenes/Moves/Ppm/Weight';
import PpmSize from 'scenes/Moves/Ppm/PPMSizeWizard';
import Review from 'scenes/Review/Review';
import Agreement from 'scenes/Legalese';

const PageNotInFlow = ({ location }) => (
  <div className="usa-grid">
    <h3>Missing Context</h3>
    You are trying to load a page that the system does not have context for.
    Please go to the home page and try again.
  </div>
);

const Placeholder = props => {
  return (
    <WizardPage
      handleSubmit={() => undefined}
      pageList={props.pageList}
      pageKey={props.pageKey}
    >
      <div className="Todo">
        <h1>Placeholder for {props.title}</h1>
        <h2>{props.description}</h2>
      </div>
    </WizardPage>
  );
};

const stub = (key, pages, description) => ({ match }) => (
  <Placeholder
    pageList={pages}
    pageKey={key}
    title={key}
    description={description}
  />
);

const always = () => true;
const hasHHG = ({ selectedMoveType }) =>
  selectedMoveType !== null && selectedMoveType !== 'PPM';
const hasPPM = ({ selectedMoveType }) =>
  selectedMoveType !== null && selectedMoveType !== 'HHG';
const isCombo = ({ selectedMoveType }) =>
  selectedMoveType !== null && selectedMoveType === 'COMBO';

const pages = {
  '/service-member/:serviceMemberId/create': {
    isInFlow: always,
    isComplete: sm =>
      sm.is_profile_complete || every([sm.rank, sm.edipi, sm.affiliation]),
    render: (key, pages) => ({ match }) => (
      <DodInfo pages={pages} pageKey={key} match={match} />
    ),
  },
  '/service-member/:serviceMemberId/name': {
    isInFlow: always,
    isComplete: sm =>
      sm.is_profile_complete || every([sm.first_name, sm.last_name]),
    render: (key, pages) => ({ match }) => (
      <SMName pages={pages} pageKey={key} match={match} />
    ),
  },
  '/service-member/:serviceMemberId/contact-info': {
    isInFlow: always,
    isComplete: sm =>
      sm.is_profile_complete ||
      (every([sm.telephone, sm.personal_email]) &&
        some([
          sm.phone_is_preferred,
          sm.email_is_preferred,
          sm.text_message_is_preferred,
        ])),
    render: (key, pages) => ({ match }) => (
      <ContactInfo pages={pages} pageKey={key} match={match} />
    ),
  },
  '/service-member/:serviceMemberId/duty-station': {
    isInFlow: always,
    isComplete: sm => sm.is_profile_complete || Boolean(sm.current_station),
    render: (key, pages) => ({ match }) => (
      <DutyStation pages={pages} pageKey={key} match={match} />
    ),
    description: 'current duty station',
  },
  '/service-member/:serviceMemberId/residence-address': {
    isInFlow: always,
    isComplete: sm => sm.is_profile_complete || Boolean(sm.residential_address),
    render: (key, pages) => ({ match }) => (
      <ResidentialAddress pages={pages} pageKey={key} match={match} />
    ),
  },
  '/service-member/:serviceMemberId/backup-mailing-address': {
    isInFlow: always,
    isComplete: sm =>
      sm.is_profile_complete || Boolean(sm.backup_mailing_address),
    render: (key, pages) => ({ match }) => (
      <BackupMailingAddress pages={pages} pageKey={key} match={match} />
    ),
  },
  '/service-member/:serviceMemberId/backup-contacts': {
    isInFlow: always,
    isComplete: sm => {
      return (
        sm.is_profile_complete || get(sm, 'backup_contacts', []).length > 0
      );
    },
    render: (key, pages) => ({ match }) => (
      <BackupContact pages={pages} pageKey={key} match={match} />
    ),
    description: 'Backup contacts',
  },
  '/service-member/:serviceMemberId/transition': {
    isInFlow: always,
    isComplete: always,
    render: (key, pages) => ({ match }) => (
      <WizardPage handleSubmit={no_op} pageList={pages} pageKey={key}>
        <TransitionToOrders />
      </WizardPage>
    ),
  },
  '/orders/': {
    isInFlow: always,
    isComplete: sm => {
      const orders = get(sm, 'orders.0', {});
      return every([
        orders.orders_type,
        orders.issue_date,
        orders.report_by_date,
        orders.new_duty_station,
      ]);
    },
    render: (key, pages) => ({ match }) => (
      <Orders pages={pages} pageKey={key} match={match} />
    ),
  },
  '/orders/upload': {
    isInFlow: always,
    isComplete: sm =>
      get(sm, 'orders.0.uploaded_orders.uploads', []).length > 0,
    render: (key, pages) => ({ match }) => (
      <UploadOrders pages={pages} pageKey={key} match={match} />
    ),
    description: 'Upload your orders',
  },
  '/orders/transition': {
    isInFlow: always,
    isComplete: always,
    render: (key, pages, description, props) => ({ match }) => (
      <WizardPage
        handleSubmit={no_op}
        isAsync={false}
        pageList={pages}
        pageKey={key}
        additionalParams={{ moveId: props.moveId }}
      >
        <TransitionToMove />
      </WizardPage>
    ),
  },
  '/moves/:moveId': {
    isInFlow: always,
    isComplete: sm => get(sm, 'orders.0.moves.0.selected_move_type', null),
    render: (key, pages) => ({ match }) => (
      <MoveType pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/schedule': {
    isInFlow: hasHHG,
    isComplete: always, //todo fix this when implemented
    render: stub,
    description: 'Pick a move date',
  },
  '/moves/:moveId/address': {
    isInFlow: hasHHG,
    isComplete: always, //todo fix this when implemented
    render: stub,
    description: 'enter your addresses',
  },

  '/moves/:moveId/ppm-transition': {
    isInFlow: isCombo,
    isComplete: always,
    render: (key, pages) => ({ match }) => (
      <WizardPage handleSubmit={no_op} pageList={pages} pageKey={key}>
        <Transition />
      </WizardPage>
    ),
  },
  '/moves/:moveId/ppm-start': {
    isInFlow: state => state.selectedMoveType === 'PPM',
    isComplete: sm => {
      const ppm = get(sm, 'orders.0.moves.0.personally_procured_moves.0', {});
      return every([
        ppm.planned_move_date,
        ppm.pickup_postal_code,
        ppm.destination_postal_code,
      ]);
    },
    render: (key, pages) => ({ match }) => (
      <PpmDateAndLocations pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/ppm-size': {
    isInFlow: hasPPM,
    isComplete: sm =>
      get(sm, 'orders.0.moves.0.personally_procured_moves.0.size', null),
    render: (key, pages) => ({ match }) => (
      <PpmSize pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/ppm-incentive': {
    isInFlow: hasPPM,
    isComplete: sm =>
      get(sm, 'orders.0.moves.0.personally_procured_moves.0.weight', null),
    render: (key, pages) => ({ match }) => (
      <PpmWeight pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/review': {
    isInFlow: always,
    isComplete: sm =>
      get(sm, 'orders.0.moves.0.status', 'DRAFT') === 'SUBMITTED',
    render: (key, pages) => ({ match }) => (
      <Review pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/agreement': {
    isInFlow: always,
    isComplete: sm =>
      get(sm, 'orders.0.moves.0.status', 'DRAFT') === 'SUBMITTED',
    render: (key, pages, description, props) => ({ match }) => {
      return <Agreement pages={pages} pageKey={key} match={match} />;
    },
  },
};
export const getPagesInFlow = state =>
  Object.keys(pages).filter(pageKey => {
    const page = pages[pageKey];
    return page.isInFlow(state);
  });

export const getNextIncompletePage = service_member => {
  const rawPath = findKey(pages, p => !p.isComplete(service_member));
  const compiledPath = generatePath(rawPath, {
    serviceMemberId: service_member.id,
    moveId: get(service_member, 'orders.0.moves.0.id', null),
  });
  return compiledPath;
};
export const getWorkflowRoutes = props => {
  const pageList = getPagesInFlow(props);
  return Object.keys(pages).map(key => {
    const currPage = pages[key];
    if (currPage.isInFlow(props)) {
      const render = currPage.render(
        key,
        pageList,
        currPage.description,
        props,
      );
      return <PrivateRoute exact path={key} key={key} render={render} />;
    } else {
      return <Route exact path={key} key={key} component={PageNotInFlow} />;
    }
  });
};
