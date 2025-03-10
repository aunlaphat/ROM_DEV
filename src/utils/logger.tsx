// Constants
const LOG_COLORS = {
    debug: '#808080',  // Gray
    info: '#1890ff',   // Primary Blue
    warn: '#faad14',   // Warning Yellow
    error: '#f5222d',  // Error Red
    success: '#52c41a' // Success Green
  };
  
  const LOG_ICONS = {
    api: '🌐',
    auth: '🔒',
    login: '🔑',
    logout: '👋',
    state: '📊',
    time: '⏱️',
    error: '❌',
    success: '✅',
    warning: '⚠️',
    navigation: '🔄',
    redux: '⚛️'  // Add Redux icon
  };
  
  // Utility functions
  const getTimestamp = () => {
    return new Date().toLocaleTimeString();
  };
  
  // Main logger object
  export const logger = {
    // Basic logging methods
    log: (level: keyof typeof LOG_COLORS, message: string, data?: any) => {
      const color = LOG_COLORS[level];
      const timestamp = getTimestamp();
      
      console.log(
        `%c${timestamp} ${message}`, 
        `color: ${color}; font-weight: bold`,
        data || ''
      );
    },
  
    // API related logs
    api: {
      request: (endpoint: string, data?: any) => {
        logger.log('info', `${LOG_ICONS.api} API Request: ${endpoint}`, data);
      },
      success: (endpoint: string, data?: any) => {
        logger.log('success', `${LOG_ICONS.success} API Success: ${endpoint}`, data);
      },
      error: (endpoint: string, error: any) => {
        logger.log('error', `${LOG_ICONS.error} API Error: ${endpoint}`, error);
      }
    },
  
    // Performance monitoring
    perf: {
      start: (label: string) => {
        console.group(`${LOG_ICONS.time} ${label}`);
        console.time(label);
      },
      end: (label: string) => {
        console.timeEnd(label);
        console.groupEnd();
      }
    },
  
    // Error handling
    error: (message: string, error?: any) => {
      logger.log('error', `${LOG_ICONS.error} ${message}`, error);
    },
  
    // State management logs
    state: {
      update: (action: string, data?: any) => {
        logger.log('info', `${LOG_ICONS.state} State Update: ${action}`, data);
      },
      error: (message: string, error?: any) => {
        logger.log('error', `${LOG_ICONS.state} State Error: ${message}`, error);
      }
    },
    
    // Redux specific logging
    redux: {
      action: (type: string, payload?: any) => {
        logger.log('info', `${LOG_ICONS.redux} Redux Action: ${type}`, payload);
      },
      state: (newState: any) => {
        logger.log('debug', `${LOG_ICONS.redux} Redux State`, newState);
      }
    },
    
    // Component logging
    component: (component: string, message: string, data?: any) => {
      logger.log('info', `[${component}] ${message}`, data);
    },
    
    // Authentication logging
    auth: {
      login: (username: string, success: boolean, data?: any) => {
        logger.log(
          success ? 'success' : 'error',
          `${LOG_ICONS.login} Auth: ${success ? 'Login Success' : 'Login Failed'} - ${username}`,
          data
        );
      },
      logout: (reason?: string) => {
        logger.log('info', `${LOG_ICONS.logout} Auth: Logout`, { reason });
      },
      check: (isAuthenticated: boolean, data?: any) => {
        logger.log(
          isAuthenticated ? 'success' : 'info',
          `${LOG_ICONS.auth} Auth Check: ${isAuthenticated ? 'Authenticated' : 'Not Authenticated'}`,
          data
        );
      }
    },
    
    // Navigation/routing logs
    navigation: {
      to: (path: string, data?: any) => {
        logger.log('info', `${LOG_ICONS.navigation} Navigation: ${path}`, data);
      }
    }
  };