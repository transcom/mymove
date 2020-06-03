import React from 'react';

import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import { ReactComponent as DocMenu } from 'shared/icon/doc-menu.svg';

const DocViewerMenu = () => (
  <div className="doc-viewer--menu doc-viewer--menu--collapsed">
    <XLightIcon />
    <DocMenu />
  </div>
);

export default DocViewerMenu;
