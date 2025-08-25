'use client';

import React, { Fragment } from 'react';
import { GlobalToaster } from '../globalToaster';

const UIProvider = ({ children }: { children: React.ReactNode }) => {
  return (
    <Fragment>
      <GlobalToaster />
      {children}
    </Fragment>
  );
};

export default UIProvider;
