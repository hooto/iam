%define app_home /opt/lessiam
%define app_user lessiam
%define app_grp  lessiam

Name:      lessiam
Version:   x.y.z
Release:   1%{?dist}
Vendor:    lessCompute.com
Summary:   lessCompute Identity Server, Enterprise PaaS Middleware
License:   Apache 2
Group:     Applications
Source0:   lessiam-x.y.z.tar.gz
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}

Requires:       redhat-lsb-core
Requires:       lessiam
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

install -m 0755 -p bin/lessiam %{buildroot}%{app_home}/bin/lessiam
install -m 0755 -p bin/lessiam-cli %{buildroot}%{app_home}/bin/lessiam-cli
install -m 0640 -p etc/lessiam.conf %{buildroot}%{app_home}/etc/lessiam.conf
install -m 0755 -p -D misc/rhel/init.d-scripts %{buildroot}%{_initrddir}/lessiam

%clean
rm -rf %{buildroot}

%pre
# Add the "lessiam" user
getent group %{app_grp} >/dev/null || groupadd -r %{app_grp}
getent passwd %{app_user} >/dev/null || \
    useradd -r -g %{app_grp} -s /sbin/nologin \
    -d %{app_home} -c "less Identity Server User"  %{app_user}

if [ $1 == 2 ]; then
    service lessiam stop
fi

%post
# Register the lessiam service
if [ $1 -eq 1 ]; then
    /sbin/chkconfig --add lessiam
    /sbin/chkconfig --level 345 lessiam on
    /sbin/chkconfig lessiam on
fi

if [ $1 -ge 1 ]; then
    service lessiam start
fi

%preun
if [ $1 = 0 ]; then
    /sbin/service lessiam stop  > /dev/null 2>&1
    /sbin/chkconfig --del lessiam
fi

%postun

%files
%defattr(-,lessiam,lessiam,-)
%dir %{app_home}
%{app_home}/
%config(noreplace) %{app_home}/etc/lessiam.conf
%{_initrddir}/lessiam
