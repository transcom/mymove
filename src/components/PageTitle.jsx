/**
 * This can be added just inside a <*Router> component to update the page title whenever the route changes.
 */

import React from 'react';

import TitleAnnouncer from './TitleAnnouncer/TitleAnnouncer';

import { useTitle } from 'hooks/custom';

export default function PageTitle() {
  useTitle();

  return <TitleAnnouncer />;
}
