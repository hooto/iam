%define app_home /opt/lessids
%define app_user lessids
%define app_grp  lessids

Name:      lessids
Version:   x.y.z
Release:   1%{?dist}
Vendor:    lessCompute.com
Summary:   lessCompute Identity Server, Enterprise PaaS Middleware
License:   Apache 2
Group:     Applications
Source0:   lessids-x.y.z.tar.gz
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}

Requires:       redhat-lsb-core
Requires:       lessids
Requires(pre):  shadow-utils
Requires(post): chkconfig

%description

%prep

%setup  -q -n %{name}-%{version}

%build
export GOROOT=/usr/local/go
export PATH=$PATH:/usr/local/go/bin
export GOPATH=~/gopath
./build

%install
rm -rf %{buildroot}

install -d %{buildroot}%{app_home}/bin
install -d %{buildroot}%{app_home}/etc
install -d %{buildroot}%{app_home}/src

cp -rp ./static %{buildroot}%{app_home}/
cp -rp ./src/views %{buildroot}%{app_home}/src/

install -m 0755 -p bin/lessids %{buildroot}%{app_home}/bin/lessids
install -m 0755 -p bin/lessids-cli %{buildroot}%{app_home}/bin/lessids-cli
install -m 0640 -p etc/lessids.conf %{buildroot}%{app_home}/etc/lessids.conf
install -m 0755 -p -D misc/rhel/init.d-scripts %{buildroot}%{_initrddir}/lessids

%clean
rm -rf %{buildroot}

%pre
# Add the "lessids" user
getent group %{app_grp} >/dev/null || groupadd -r %{app_grp}
getent passwd %{app_user} >/dev/null || \
    useradd -r -g %{app_grp} -s /sbin/nologin \
    -d %{app_home} -c "less Identity Server User"  %{app_user}

if [ $1 == 2 ]; then
    service lessids stop
fi

%post
# Register the lessids service
if [ $1 -eq 1 ]; then
    /sbin/chkconfig --add lessids
    /sbin/chkconfig --level 345 lessids on
    /sbin/chkconfig lessids on
fi

if [ $1 -ge 1 ]; then
    service lessids start
fi

%preun
if [ $1 = 0 ]; then
    /sbin/service lessids stop  > /dev/null 2>&1
    /sbin/chkconfig --del lessids
fi

%postun

%files
%defattr(-,lessids,lessids,-)
%dir %{app_home}
%{app_home}/
%config(noreplace) %{app_home}/etc/lessids.conf
%{_initrddir}/lessids
